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

type Player struct {
	guildID     string
	connector   MusicConnector
	hub         *messaging.Hub
	botAudio    *bot.BotAudio
	playlist    []*Track
	progress    time.Duration
	stateChange chan MusicPlayerState
	currentUrl  string
	stopped     bool
	preventSkip bool
	logger log.Logger
}

type PlayerManager struct {
	players   map[string]*Player
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

var ErrEmptyPlaylist = errors.New("empty playlist")

func NewMusicPlayerManager(connector MusicConnector, hub *messaging.Hub, bot *bot.Bot) *PlayerManager {
	return &PlayerManager{
		players:   make(map[string]*Player),
		connector: connector,
		hub:       hub,
		bot:       bot,
	}
}

func (m *PlayerManager) GetPlayer(guildID string) (*Player, error) {
	if guildID == "" {
		return nil, errors.New("guild id cannot be empty")
	}

	for _, p := range m.players {
		if p.guildID == guildID {
			return p, nil
		}
	}

	p := &Player{
		guildID:     guildID,
		connector:   m.connector,
		hub:         m.hub,
		botAudio:    m.bot.GetGuildVoice(guildID),
		stateChange: make(chan MusicPlayerState),
		playlist:    []*Track{},
		stopped:     false,
		preventSkip: false,
	}
	m.players[guildID] = p
	p.subscribeToAudioStatusChange()
	return p, nil
}

func (m *PlayerManager) Search(query string, page uint) (*MusicSearchResult, error) {
	return m.connector.Search(query, page)
}

func (p *Player) AddToPlaylist(ID string, elemType MusicElementType) error {
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

func (p *Player) MoveInPlaylist(from, to int) {
	el := p.playlist[from]
	arr := append(p.playlist[:from], p.playlist[from+1:]...)
	arr = append(arr[:to+1], arr[to:]...)
	arr[to] = el
	p.playlist = arr
	p.propagateState()
}

func (p *Player) RemoveFromPlaylist(i int) {
	p.playlist = append(p.playlist[:i], p.playlist[i+1:]...)
	p.propagateState()
}

func (p *Player) Play() error {

	if len(p.playlist) <= 0 {
		p.logger.Debug()
		return ErrEmptyPlaylist
	}

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
	var err error
	p.currentUrl, err = p.connector.GetStreamURL(p.playlist[0].ID)
	if err != nil {
		return err
	}

	// Start the stream
	end, progress, err := p.botAudio.Play(p.currentUrl, 0)
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
		if p.stopped || p.preventSkip {
			p.preventSkip = false
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

func (p *Player) Pause() {
	p.botAudio.Pause()
}

func (p *Player) SetTime(t time.Duration) error {
	p.preventSkip = true
	// Start the stream
	end, progress, err := p.botAudio.Play(p.currentUrl, t)
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
		if p.stopped || p.preventSkip {
			p.preventSkip = false
			return
		}
		p.playlist = p.playlist[1:]
		if len(p.playlist) > 0 {
			p.Play()
		}
	}()
	return nil
}

func (p *Player) Stop() {
	p.stopped = true
	p.botAudio.Stop()
}

func (p *Player) Skip() error {
	if len(p.playlist) < 2 {
		return nil
	}

	p.stopped = false
	p.botAudio.Stop()
	return nil
}

func (p *Player) ClearPlaylist() {
	p.playlist = []*Track{}
	p.propagateState()
}

func (p *Player) Playlist() []*Track {
	return p.playlist
}

func (p *Player) BotAudio() *bot.BotAudio {
	return p.botAudio
}

func (p *Player) SubscribeToPlayerState() chan MusicPlayerState {
	return p.stateChange
}

func (p *Player) propagateState() {
	payload := MusicPlayerState{
		Status:   p.botAudio.Status(),
		Playlist: p.playlist,
		Progress: p.progress.Milliseconds(),
	}

	go p.updatePlayerState(payload)

	m := messaging.Message{
		MessageType: "musicPlayer/updatePlayerState",
		Payload:     payload,
		RoomID:      p.guildID,
	}
	p.hub.Send(m)
}

func (p *Player) updatePlayerState(state MusicPlayerState) {
	p.stateChange <- state
}

func (p *Player) subscribeToAudioStatusChange() {
	go func() {
		for {
			<-p.botAudio.StatusChange()
			p.propagateState()
		}
	}()
}
