package bot

import (
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/config"
	"github.com/ggresillion/discordsoundboard/backend/internal/messaging"
)

type Bot struct {
	hub              *messaging.Hub
	session          *discordgo.Session
	voiceStates      map[string][]*discordgo.VoiceState
	voiceConnections map[string]*discordgo.VoiceConnection
}

func NewBot(hub *messaging.Hub) *Bot {
	return &Bot{hub: hub, session: nil, voiceStates: make(map[string][]*discordgo.VoiceState), voiceConnections: make(map[string]*discordgo.VoiceConnection)}
}

func (b *Bot) StartBot() error {

	token := config.Get().BotToken

	if token == "" {
		return fmt.Errorf("no token provided")
	}

	var err error
	b.session, err = discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	b.session.AddHandler(ready)
	b.session.AddHandler(b.guildCreate)

	b.session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuildMembers

	err = b.session.Open()
	if err != nil {
		return fmt.Errorf("error opening Discord session: %w", err)
	}

	b.registerCommands()

	log.Println("bot is now running")
	return nil
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	s.UpdateGameStatus(0, "!dsb")
}

func (b *Bot) guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	b.voiceStates[event.Guild.ID] = event.VoiceStates
}

func (b *Bot) JoinUserChannel(guildID, userID string) (*discordgo.VoiceConnection, error) {
	g, err := b.session.State.Guild(guildID)
	if err != nil {
		return nil, err
	}
	for _, v := range g.VoiceStates {
		if v.UserID == userID {
			defer log.Printf("joinned voice channel for guildID %s and userID %s", guildID, userID)
			vc, err := b.session.ChannelVoiceJoin(guildID, v.ChannelID, false, true)
			if err != nil {
				return nil, err
			}
			b.voiceConnections[guildID] = vc
			return vc, nil
		}
	}
	return nil, errors.New("user not in a voice channel")
}

func (b *Bot) GetGuild(id string) (*discordgo.Guild, error) {
	return b.session.Guild(id)
}

func (b *Bot) GetVoiceConnectionOrJoin(guildID, userID string) (*discordgo.VoiceConnection, error) {
	vc := b.voiceConnections[guildID]
	if vc != nil {
		return vc, nil
	}
	return b.JoinUserChannel(userID, guildID)
}

func (b *Bot) GetVoiceConnection(guildID string) (*discordgo.VoiceConnection, error) {
	vc := b.voiceConnections[guildID]
	if vc != nil {
		return vc, nil
	}
	return nil, fmt.Errorf("bot is not connected in a voice channel for guildID %s", guildID)
}
