package music

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal"
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
	hub       *internal.Hub
}

type MusicPlayerManager struct {
	players   map[string]*MusicPlayer
	connector MusicConnector
	hub       *internal.Hub
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
}

func NewMusicPlayerManager(connector MusicConnector, hub *internal.Hub) *MusicPlayerManager {
	return &MusicPlayerManager{
		players:   make(map[string]*MusicPlayer),
		connector: connector,
		hub:       hub,
	}
}

func (m *MusicPlayerManager) GetPlayer(guildID string) *MusicPlayer {
	for _, p := range m.players {
		if p.guildID == guildID {
			return p
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
	}
	m.players[guildID] = p
	return p
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

func (p *MusicPlayer) AddToPlaylist(ID string) error {
	t, err := p.connector.Find(ID)
	if err != nil {
		return err
	}

	p.state.Playlist = append(p.state.Playlist, *t)
	p.PropagateState()
	return nil
}

func (p *MusicPlayer) RemoveFromPlaylist(id string) {
	for i, t := range p.state.Playlist {
		if t.ID == id {
			p.state.Playlist = removeIndex(p.state.Playlist, i)
		}
	}
	p.PropagateState()
}

func (p *MusicPlayer) PropagateState() {
	m := internal.Message{
		MessageType: "musicPlayer/updateState",
		Payload:     p.state,
		RoomID:      p.guildID,
	}
	p.hub.Send(m)
}

func removeIndex(s []Track, index int) []Track {
	return append(s[:index], s[index+1:]...)
}
