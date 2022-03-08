package api

import (
	"log"
	"net/http"

	"github.com/ggresillion/discordsoundboard/backend/internal"
	"github.com/gorilla/websocket"
)

type WSAPI struct {
	hub *internal.Hub
}

func NewWSAPI(hub *internal.Hub) *WSAPI {
	return &WSAPI{hub: hub}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (a *WSAPI) writePump(conn *websocket.Conn, client *internal.Client) {
	s := client.ToSend()
	for {
		m := <-s
		conn.WriteJSON(m)
	}
}

func (a *WSAPI) readPump(conn *websocket.Conn, client *internal.Client) {
	for {
		m := internal.Message{}
		err := conn.ReadJSON(&m)
		if err != nil {
			conn.Close()
			log.Println(err)
			return
		}

		client.Receive(m)
	}
}

func (a *WSAPI) handleConnection(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := internal.NewClient(a.hub)
	a.hub.Register(client)

	go a.readPump(conn, client)
	go a.writePump(conn, client)
}
