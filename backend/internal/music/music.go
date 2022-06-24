package music

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/messaging"
)

type Connector interface {
	Search(q string, p uint) (*SearchResult, error)
	FindTrack(ID string) (*Track, error)
	FindPlaylistTracks(ID string) ([]*Track, error)
	FindAlbumTracks(ID string) ([]*Track, error)
	FindArtistTopTracks(ID string) ([]*Track, error)
	GetStreamURL(ID string) (string, error)
}

type Player struct {
	guildID        string
	connector      Connector
	hub            *messaging.Hub
	botAudio       *bot.Audio
	playlist       []*Track
	progress       time.Duration
	progressChange chan time.Duration
	stateChange    chan PlayerState
	currentUrl     string
	stopped        bool
	preventSkip    bool
}

type PlayerManager struct {
	players   map[string]*Player
	connector Connector
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

type SearchResult struct {
	Tracks    []Track    `json:"tracks"`
	Playlists []Playlist `json:"playlists"`
	Albums    []Album    `json:"albums"`
	Artists   []Artist   `json:"artists"`
}

type PlayerState struct {
	Status   bot.AudioStatus `json:"status"`
	Playlist []*Track        `json:"playlist"`
}

type ElementType int

const (
	TrackElement ElementType = iota
	ArtistElement
	AlbumElement
	PlaylistElement
)

var ErrEmptyPlaylist = errors.New("empty playlist")

func NewPlayerManager(connector Connector, hub *messaging.Hub, bot *bot.Bot) *PlayerManager {
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
		guildID:        guildID,
		connector:      m.connector,
		hub:            m.hub,
		botAudio:       m.bot.GetGuildVoice(guildID),
		stateChange:    make(chan PlayerState),
		playlist:       []*Track{},
		stopped:        false,
		preventSkip:    false,
		progressChange: make(chan time.Duration),
	}
	m.players[guildID] = p
	p.subscribeToAudioStatusChange()
	return p, nil
}

func (m *PlayerManager) Search(query string, page uint) (*SearchResult, error) {
	return m.connector.Search(query, page)
}

func (p *Player) AddToPlaylist(ID string, elemType ElementType) error {
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
		log.Debugf("[%s] (play) playlist is empty", p.guildID)
		return ErrEmptyPlaylist
	}

	switch p.botAudio.Status() {
	case bot.Paused:
		log.Debugf("[%s] (play) unpausing stream", p.guildID)
		p.botAudio.UnPause()
		return nil
	case bot.Playing:
		log.Debugf("[%s] (play) already playing.", p.guildID)
		return nil
	}

	// Make sure the stream is not stopped
	p.stopped = false

	// Get the stream url from source
	var err error
	p.currentUrl, err = p.connector.GetStreamURL(p.playlist[0].ID)
	if err != nil {
		log.Errorf("[%s] (play) error getting audio stream url: %e", p.guildID, err)
		return err
	}
	log.Debugf("[%s] (play) got audio stream url: %s", p.guildID, p.currentUrl)

	// Start the stream
	stream, err := p.botAudio.Play(p.currentUrl, 0)
	if err != nil {
		log.Errorf("[%s] error starting audio stream: %e", p.guildID, err)
		return err
	}
	log.Debugf("[%s] (play) playing stream", p.guildID)

	p.handleProgress(stream)
	p.handleStreamEnd(stream)
	p.propagateState()
	return nil
}

// handleProgress
// handles the stream progress change
func (p *Player) handleProgress(stream *bot.Stream) {
	go func() {
		for {
			pr, more := <-stream.Progress
			if !more {
				return
			}
			p.progress = pr
			p.progressChange <- pr
		}
	}()
}

// handleStreamEnd
// handles the stream ending
func (p *Player) handleStreamEnd(stream *bot.Stream) {
	go func() {
		err := <-stream.End
		if err == nil {
			log.Debugf("[%s] stream ended", p.guildID)
		} else {
			log.Debugf("[%s] stream ended with error: %e", p.guildID, err)
		}

		// Do nothing if the player is stopped
		if p.stopped || p.preventSkip {
			log.Debugf("[%s] stopping player", p.guildID)
			p.preventSkip = false
			return
		}
		p.playlist = p.playlist[1:]
		if len(p.playlist) > 0 {
			log.Debugf("[%s] playing next song", p.guildID)
			err := p.Play()
			if err != nil {
				log.Errorf("[%s] error playing next song: %e", p.guildID, err)
			}
		}
	}()
}

func (p *Player) Pause() {
	log.Debugf("[%s] (pause) pausing stream", p.guildID)
	p.botAudio.Pause()
}

func (p *Player) SetTime(t time.Duration) error {
	p.preventSkip = true
	// Start the stream
	stream, err := p.botAudio.Play(p.currentUrl, t)
	if err != nil {
		return fmt.Errorf("error starting audio streaming: %v", err)
	}
	p.handleProgress(stream)
	p.handleStreamEnd(stream)
	return nil
}

func (p *Player) Stop() {
	p.stopped = true
	p.botAudio.Stop()
}

func (p *Player) Skip() error {
	if len(p.playlist) < 2 {
		log.Debugf("[%s] (skip) could not skip, no next music", p.guildID)
		return nil
	}

	p.stopped = false
	p.botAudio.Stop()
	log.Debugf("[%s] (skip) skipped next music", p.guildID)
	return nil
}

func (p *Player) ClearPlaylist() {
	p.playlist = []*Track{}
	p.propagateState()
}

func (p *Player) Playlist() []*Track {
	return p.playlist
}

func (p *Player) BotAudio() *bot.Audio {
	return p.botAudio
}

func (p *Player) SubscribeToPlayerState() chan PlayerState {
	return p.stateChange
}

func (p *Player) Progress() time.Duration {
	return p.progress
}

func (p *Player) SubscribeToProgress() chan time.Duration {
	c := make(chan time.Duration)
	go func() {
		for {
			val := <-p.progressChange
			c <- val
		}
	}()
	return c
}

func (p *Player) propagateState() {
	state := PlayerState{
		Status:   p.botAudio.Status(),
		Playlist: p.playlist,
	}
	log.Debugf("[%s] (propagate) next state: {status: %s, playlist: %d element(s)}", p.guildID, state.Status,
		len(state.Playlist))
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
