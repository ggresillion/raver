package music

import (
	"errors"
	"fmt"
	"log"

	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/messaging"
)

type MusicConnector interface {
	Search(q string, p uint) ([]Track, error)
	Find(ID string) (*Track, error)
	GetStreamURL(query string) (string, error)
}

type MusicPlayer struct {
	guildID   string
	connector MusicConnector
	hub       *messaging.Hub
	botAudio  *bot.BotAudio
	playlist  []Track
}

type MusicPlayerManager struct {
	players   map[string]*MusicPlayer
	connector MusicConnector
	hub       *messaging.Hub
	bot       *bot.Bot
}

type Track struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	Thumbnail string `json:"thumbnail"`
	Duration  uint   `json:"duration"`
	URL       string `json:"url"`
}

type MusicPlayerState struct {
	Status   bot.AudioStatus `json:"status"`
	Playlist []Track         `json:"playlist"`
}

func NewMusicPlayerManager(connector MusicConnector, hub *messaging.Hub, bot *bot.Bot) *MusicPlayerManager {
	return &MusicPlayerManager{
		players:   make(map[string]*MusicPlayer),
		connector: connector,
		hub:       hub,
		bot:       bot,
	}
}

func (m *MusicPlayerManager) GetPlayer(guildID string) (*MusicPlayer, error) {
	if guildID == "" {
		return nil, errors.New("guild id cannot be empty")
	}

	for _, p := range m.players {
		if p.guildID == guildID {
			return p, nil
		}
	}

	p := &MusicPlayer{
		guildID:   guildID,
		connector: m.connector,
		hub:       m.hub,
		botAudio:  m.bot.GetGuildVoice(guildID),
		playlist:  []Track{},
	}
	m.players[guildID] = p
	p.subscribeToAudioStatusChange()
	return p, nil
}

func (m *MusicPlayerManager) Search(query string, page uint) ([]Track, error) {
	return m.connector.Search(query, page)
}

func (p *MusicPlayer) AddToPlaylist(ID string) (*Track, error) {
	t, err := p.connector.Find(ID)
	if err != nil {
		return nil, err
	}

	p.playlist = append(p.playlist, *t)
	p.propagateState()
	return t, nil
}

func (p *MusicPlayer) MoveInPlaylist(from, to int) {
	p.playlist = move(p.playlist, from, to)
	p.propagateState()
}

func (p *MusicPlayer) RemoveFromPlaylist(ID string) {
	for i, t := range p.playlist {
		if t.ID == ID {
			p.playlist = remove(p.playlist, i)
		}
	}
	p.propagateState()
}

func (p *MusicPlayer) Play() error {

	switch p.botAudio.Status() {
	case bot.Paused:
		p.botAudio.UnPause()
		return nil
	case bot.Playing:
		return nil
	}

	// Get the stream url from source
	url, err := p.connector.GetStreamURL(p.playlist[0].ID)
	if err != nil {
		return err
	}

	// Start the stream
	end, err := p.botAudio.Play(url)
	if err != nil {
		return fmt.Errorf("error starting audio streaming: %v", err)
	}

	// Listen for stream end, skiping to next song if any
	go func() {
		err := <-end
		if err != nil {
			log.Printf("error during audio streaming: %v", err)
		}
		p.playlist = p.playlist[1:]
		if len(p.playlist) > 0 {
			p.Play()
		}
	}()

	p.propagateState()
	return nil
}

func (p *MusicPlayer) Pause() {
	p.botAudio.Pause()
}

func (p *MusicPlayer) Stop() {
	p.botAudio.Stop()
}

func (p *MusicPlayer) Skip() error {
	if len(p.playlist) < 2 {
		return nil
	}

	p.botAudio.Stop()
	return nil
}

func (p *MusicPlayer) ClearPlaylist() {
	p.playlist = []Track{}
	p.propagateState()
}

func (p *MusicPlayer) Playlist() []Track {
	return p.playlist
}

func (p *MusicPlayer) BotAudio() *bot.BotAudio {
	return p.botAudio
}

func (p *MusicPlayer) propagateState() {
	payload := MusicPlayerState{
		Status:   p.botAudio.Status(),
		Playlist: p.playlist,
	}

	m := messaging.Message{
		MessageType: "musicPlayer/updatePlayerState",
		Payload:     payload,
		RoomID:      p.guildID,
	}
	p.hub.Send(m)
}

func (p *MusicPlayer) subscribeToAudioStatusChange() {
	go func() {
		for {
			<-p.botAudio.StatusChange()
			p.propagateState()
		}
	}()
}

func insert(array []Track, value Track, index int) []Track {
	return append(array[:index], append([]Track{value}, array[index:]...)...)
}

func remove(array []Track, index int) []Track {
	return append(array[:index], array[index+1:]...)
}

func move(array []Track, srcIndex int, dstIndex int) []Track {
	value := array[srcIndex]
	return insert(remove(array, srcIndex), value, dstIndex)
}
