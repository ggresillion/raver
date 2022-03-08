package bot

import (
	"errors"
	"fmt"

	"github.com/ggresillion/discordsoundboard/backend/internal"
	"github.com/ggresillion/discordsoundboard/backend/internal/discord"
)

type BotSubscriber struct {
	bot *Bot
}

func NewBotSubscriber(bot *Bot) *BotSubscriber {
	return &BotSubscriber{bot}
}

func (s *BotSubscriber) SubscribeToIncomingMessages() {
	c := s.bot.hub.Subscribe()
	go func() {
		for {
			m := <-c
			switch m.MessageType {
			case "websocket/join":
				s.HandleJoinMessage(m)
			}
		}
	}()
}

func (s *BotSubscriber) HandleJoinMessage(m internal.Message) {
	p := &internal.JoinRoomPayload{}
	m.CastPayload(p)

	if p.RoomID == "" {
		m.Error(errors.New("you must specify a room id"))
		return
	}

	dc := discord.NewDiscordClient(m.Token)
	guilds, err := dc.GetGuilds()
	if err != nil {
		m.Error(err)
		return
	}

	for _, guild := range guilds {
		if guild.ID == p.RoomID {
			m.Client.Join(p.RoomID)
			m.Ok()
			return
		}
	}
	m.Error(errors.New(fmt.Sprint("could not find guild with id ", p.RoomID)))
}
