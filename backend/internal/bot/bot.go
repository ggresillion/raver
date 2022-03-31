package bot

import (
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
	ready            bool
}

func NewBot(hub *messaging.Hub) *Bot {
	return &Bot{hub: hub, session: nil, voiceStates: make(map[string][]*discordgo.VoiceState), voiceConnections: make(map[string]*discordgo.VoiceConnection), ready: false}
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
	b.ready = true

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
