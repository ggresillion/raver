package music

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/messaging"
)

type MusicConnector interface {
	Search(q string, p uint) ([]Track, error)
	Find(ID string) (*Track, error)
	Play(v *discordgo.VoiceConnection, id string) error
}

type MusicPlayer struct {
	guildID   string
	connector MusicConnector
	state     *MusicPlayerState
	hub       *messaging.Hub
	bot       *bot.Bot
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
}

type MusicPlayerState struct {
	Status   MusicPlayerStatus `json:"status"`
	Playlist []Track           `json:"playlist"`
	Progress int               `json:"progress"`
}

func NewMusicPlayerManager(connector MusicConnector, hub *messaging.Hub, bot *bot.Bot) *MusicPlayerManager {
	m := &MusicPlayerManager{
		players:   make(map[string]*MusicPlayer),
		connector: connector,
		hub:       hub,
		bot:       bot,
	}
	m.registerCommands()
	return m
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
		state: &MusicPlayerState{
			Status:   IDLE,
			Playlist: make([]Track, 0),
		},
		bot: m.bot,
	}
	m.players[guildID] = p
	return p, nil
}

func (m *MusicPlayerManager) Search(query string, page uint) ([]Track, error) {
	return m.connector.Search(query, page)
}

func (p *MusicPlayer) GetState() *MusicPlayerState {
	return p.state
}

func (p *MusicPlayer) SetState(newState *MusicPlayerState) *MusicPlayerState {
	p.state = newState
	p.PropagateState()
	return newState
}

func (p *MusicPlayer) AddToPlaylist(ID string) (*Track, error) {
	t, err := p.connector.Find(ID)
	if err != nil {
		return nil, err
	}

	p.state.Playlist = append(p.state.Playlist, *t)
	p.PropagateState()
	return t, nil
}

func (p *MusicPlayer) MoveInPlaylist(from, to int) {
	p.state.Playlist = move(p.state.Playlist, from, to)
	p.PropagateState()
}

func (p *MusicPlayer) RemoveFromPlaylist(ID string) {
	for i, t := range p.state.Playlist {
		if t.ID == ID {
			p.state.Playlist = remove(p.state.Playlist, i)
		}
	}
	p.PropagateState()
}

func (p *MusicPlayer) Play() error {
	s := p.GetState()
	vc, err := p.bot.GetVoiceConnection(p.guildID)
	if err != nil {
		return err
	}

	err = p.connector.Play(vc, p.state.Playlist[0].ID)
	if err != nil {
		return err
	}
	s.Status = Playing
	p.SetState(s)
	return nil
}

func (p *MusicPlayer) Pause() {
	s := p.GetState()
	s.Status = Paused
	p.SetState(s)
}

func (p *MusicPlayer) PropagateState() {
	m := messaging.Message{
		MessageType: "musicPlayer/updatePlayerState",
		Payload:     p.state,
		RoomID:      p.guildID,
	}
	p.hub.Send(m)
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
