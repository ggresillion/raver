package music

//import "github.com/ggresillion/discordsoundboard/backend/internal/messaging"
//
//type MusicGateway struct {
//	manager *PlayerManager
//	hub     *messaging.Hub
//}
//
//func NewMusicSubscriber(manager *PlayerManager, hub *messaging.Hub) *MusicGateway {
//	return &MusicGateway{manager, hub}
//}
//
//func (g *MusicGateway) SubscribeToIncomingMessages() {
//	c := g.hub.Subscribe()
//	go func() {
//		for {
//			m := <-c
//			switch m.MessageType {
//			case "musicPlayer/updatePlayerState":
//				g.HandleStateUpdate(m)
//			}
//		}
//	}()
//}
//
//func (g *MusicGateway) HandleStateUpdate(m messaging.Message) {
//	state := &PlayerState{}
//	m.CastPayload(state)
//
//	player, err := g.manager.GetPlayer(m.RoomID)
//	if err != nil {
//		m.Error(err)
//	}
//
//	player.playlist = state.Playlist
//	player.propagateState()
//}
