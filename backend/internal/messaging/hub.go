package messaging

import "encoding/json"

type MessageType int

type Client struct {
	send     chan Message
	received chan Message
	hub      *Hub
}

type Room struct {
	id          string
	clients     []*Client
	subscribers []chan Message
}

type Message struct {
	ID          string      `json:"id"`
	MessageType string      `json:"type"`
	Token       string      `json:"token"`
	Payload     interface{} `json:"payload"`
	RoomID      string      `json:"room"`
	Client      *Client     `json:"-"`
}

type JoinRoomPayload struct {
	RoomID string `json:"room"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}

type Hub struct {
	clients     []*Client
	rooms       []*Room
	subscribers []chan Message
}

func NewHub() *Hub {
	return &Hub{}
}

func NewClient(hub *Hub) *Client {
	return &Client{
		send:     make(chan Message),
		received: make(chan Message),
		hub:      hub,
	}
}

func (h *Hub) Register(c *Client) {
	h.clients = append(h.clients, c)
}

func (h *Hub) Subscribe() chan Message {
	c := make(chan Message)
	h.subscribers = append(h.subscribers, c)
	return c
}

func (h *Hub) SubscribeToRoom(roomID string) chan Message {
	c := make(chan Message)
	found := false
	for _, r := range h.rooms {
		if r.id == roomID {
			r.subscribers = append(r.subscribers, c)
			found = true
		}
	}
	if !found {
		r := &Room{}
		r.subscribers = append(r.subscribers, c)
		h.rooms = append(h.rooms, r)
	}
	return c
}

func (h *Hub) Broadcast(m Message) {
	for _, c := range h.clients {
		c.Send(m)
	}
}

func (h *Hub) Send(m Message) {
	if m.RoomID == "" {
		h.Broadcast(m)
	}
	for _, r := range h.rooms {
		if r.id == m.RoomID {
			for _, c := range r.clients {
				c.Send(m)
			}
		}
	}
}

func (c *Client) Send(m Message) {
	c.send <- m
}

func (c *Client) Receive(m Message) {
	m.Client = c

	if m.RoomID != "" {
		for _, r := range c.hub.rooms {
			if r.id == m.RoomID {
				for _, c := range r.subscribers {
					c <- m
				}
			}
		}
		return
	}
	for _, c := range c.hub.subscribers {
		c <- m
	}
}

func (c *Client) ToSend() chan Message {
	return c.send
}

func (c *Client) Join(id string) {
	found := false
	for _, r := range c.hub.rooms {
		if r.id == id {
			r.clients = append(r.clients, c)
			found = true
		}
	}
	if !found {
		r := &Room{id: id}
		c.hub.rooms = append(c.hub.rooms, r)
		r.clients = append(r.clients, c)
	}
}

func (m *Message) CastPayload(t interface{}) error {
	p, err := json.Marshal(m.Payload)
	if err != nil {
		return err
	}
	json.Unmarshal(p, t)
	return nil
}

func (m *Message) Reply(r Message) {
	m.Client.Send(r)
}

func (m *Message) Ok() {
	m.Client.Send(Message{
		ID:          m.ID,
		MessageType: "websocket/ok",
	})
}

func (m *Message) Error(err error) {
	m.Client.Send(Message{
		ID:          m.ID,
		MessageType: "websocket/error",
		Payload: &ErrorPayload{
			Message: err.Error(),
		},
	})
}
