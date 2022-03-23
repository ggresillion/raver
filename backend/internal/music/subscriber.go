package music

import "github.com/ggresillion/discordsoundboard/backend/internal/messaging"

type MusicSubscriber struct {
	manager *MusicPlayerManager
	hub     *messaging.Hub
}

func NewMusicSubscriber(manager *MusicPlayerManager, hub *messaging.Hub) *MusicSubscriber {
	return &MusicSubscriber{manager, hub}
}

func (s *MusicSubscriber) SubscribeToIncomingMessages() {
	c := s.hub.Subscribe()
	go func() {
		for {
			m := <-c
			switch m.MessageType {
			case "musicPlayer/updatePlayerState":
				s.HandleStateUpdate(m)
			}
		}
	}()
}

func (s *MusicSubscriber) HandleStateUpdate(m messaging.Message) {
	state := &MusicPlayerState{}
	m.CastPayload(state)

	player, err := s.manager.GetPlayer(m.RoomID)
	if err != nil {
		m.Error(err)
	}

	player.state = state
	player.PropagateState()
}
