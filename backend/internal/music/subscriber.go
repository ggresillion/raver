package music

import (
	"github.com/ggresillion/discordsoundboard/backend/internal"
)

type MusicSubscriber struct {
	manager *MusicPlayerManager
	hub     *internal.Hub
}

func NewMusicSubscriber(manager *MusicPlayerManager, hub *internal.Hub) *MusicSubscriber {
	return &MusicSubscriber{manager, hub}
}

func (s *MusicSubscriber) SubscribeToIncomingMessages() {
	c := s.hub.Subscribe()
	go func() {
		for {
			m := <-c
			switch m.MessageType {
			case "musicPlayer/updateState":
				s.HandleStateUpdate(m)
			}
		}
	}()
}

func (s *MusicSubscriber) HandleStateUpdate(m internal.Message) {
	state := &MusicPlayerState{}
	m.CastPayload(state)

	player := s.manager.GetPlayer(m.RoomID)

	player.state = state
	player.PropagateState()
}
