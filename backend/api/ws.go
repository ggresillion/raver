package api

import (
	"log"
	"net/http"

	"github.com/ggresillion/discordsoundboard/backend/internal/messaging"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WSAPI struct {
	hub *messaging.Hub
}

func NewWSAPI(hub *messaging.Hub) *WSAPI {
	return &WSAPI{hub: hub}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (a *WSAPI) writePump(conn *websocket.Conn, client *messaging.Client) {
	s := client.ToSend()
	for {
		m := <-s
		conn.WriteJSON(m)
	}
}

func (a *WSAPI) readPump(conn *websocket.Conn, client *messaging.Client) {
	for {
		m := messaging.Message{}
		err := conn.ReadJSON(&m)
		if err != nil {
			conn.Close()
			log.Println(err)
			return
		}

		client.Receive(m)
	}
}

func (a *WSAPI) HandleConnection(c echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := messaging.NewClient(a.hub)
	a.hub.Register(client)

	go a.readPump(conn, client)
	go a.writePump(conn, client)
	return nil
}
