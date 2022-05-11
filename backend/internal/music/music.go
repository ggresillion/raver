package music

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/messaging"
)

type MusicConnector interface {
	Search(q string, p uint) (*MusicSearchResult, error)
	FindTrack(ID string) (*Track, error)
	FindPlaylistTracks(ID string) ([]*Track, error)
	FindAlbumTracks(ID string) ([]*Track, error)
	FindArtistTopTracks(ID string) ([]*Track, error)
	GetStreamURL(ID string) (string, error)
}

type MusicPlayer struct {
	guildID   string
	connector MusicConnector
	hub       *messaging.Hub
	botAudio  *bot.BotAudio
	playlist  []*Track
	progress  time.Duration
	stopped   bool
}

type MusicPlayerManager struct {
	players   map[string]*MusicPlayer
	connector MusicConnector
	hub       *messaging.Hub
	bot       *bot.Bot
}

type Track struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Thumbnail string   `json:"thumbnail"`
	Artists   []Artist `json:"artists"`
	Album     Album    `json:"album"`
	Duration  uint     `json:"duration"`
}

type Playlist struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
}
type Album struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Thumbnail string   `json:"thumbnail"`
	Artists   []Artist `json:"artists"`
}
type Artist struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
}

type MusicSearchResult struct {
	Tracks    []Track    `json:"tracks"`
	Playlists []Playlist `json:"playlists"`
	Albums    []Album    `json:"albums"`
	Artists   []Artist   `json:"artists"`
}

type MusicPlayerState struct {
	Status   bot.AudioStatus `json:"status"`
	Playlist []*Track        `json:"playlist"`
	Progress int64           `json:"progress"`
}

type MusicElementType int

const (
	TrackElement MusicElementType = iota
	ArtistElement
	AlbumElement
	PlaylistElement
)

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
		playlist:  []*Track{},
		stopped:   false,
	}
	m.players[guildID] = p
	p.subscribeToAudioStatusChange()
	return p, nil
}

func (m *MusicPlayerManager) Search(query string, page uint) (*MusicSearchResult, error) {
	return m.connector.Search(query, page)
}

func (p *MusicPlayer) AddToPlaylist(ID string, elemType MusicElementType) error {
	switch elemType {
	case TrackElement:
		track, err := p.connector.FindTrack(ID)
		if err != nil {
			return err
		}
		p.playlist = append(p.playlist, track)
	case PlaylistElement:
		tracks, err := p.connector.FindPlaylistTracks(ID)
		if err != nil {
			return err
		}
		p.playlist = tracks
	case AlbumElement:
		tracks, err := p.connector.FindAlbumTracks(ID)
		if err != nil {
			return err
		}
		p.playlist = tracks
	case ArtistElement:
		tracks, err := p.connector.FindArtistTopTracks(ID)
		if err != nil {
			return err
		}
		p.playlist = tracks
	}

	p.propagateState()
	return nil
}

func (p *MusicPlayer) MoveInPlaylist(from, to int) {
	el := p.playlist[from]
	arr := append(p.playlist[:from], p.playlist[from+1:]...)
	arr = append(arr[:to+1], arr[to:]...)
	arr[to] = el
	p.playlist = arr
	p.propagateState()
}

func (p *MusicPlayer) RemoveFromPlaylist(i int) {
	p.playlist = append(p.playlist[:i], p.playlist[i+1:]...)
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

	// Make sure the stream is not stopped
	p.stopped = false

	// Get the stream url from source
	url, err := p.connector.GetStreamURL(p.playlist[0].ID)
	if err != nil {
		return err
	}

	// Start the stream
	end, progress, err := p.botAudio.Play(url)
	if err != nil {
		return fmt.Errorf("error starting audio streaming: %v", err)
	}

	// Updates the stream progress
	go func() {
		for {
			progress, more := <-progress
			if !more {
				return
			}
			p.progress = progress
			p.propagateState()
		}
	}()

	// Listen for stream end, skiping to next song if any
	go func() {
		err := <-end
		if err != nil {
			log.Printf("error during audio streaming: %v", err)
		}

		// Do nothing if the player is stopped
		if p.stopped {
			return
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

func (p *MusicPlayer) SetTime(t time.Duration) {
	p.stopped = true
	p.botAudio.SetTime(t)
}

func (p *MusicPlayer) Stop() {
	p.stopped = true
	p.botAudio.Stop()
}

func (p *MusicPlayer) Skip() error {
	if len(p.playlist) < 2 {
		return nil
	}

	p.stopped = false
	p.botAudio.Stop()
	return nil
}

func (p *MusicPlayer) ClearPlaylist() {
	p.playlist = []*Track{}
	p.propagateState()
}

func (p *MusicPlayer) Playlist() []*Track {
	return p.playlist
}

func (p *MusicPlayer) BotAudio() *bot.BotAudio {
	return p.botAudio
}

func (p *MusicPlayer) propagateState() {
	payload := MusicPlayerState{
		Status:   p.botAudio.Status(),
		Playlist: p.playlist,
		Progress: p.progress.Milliseconds(),
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
