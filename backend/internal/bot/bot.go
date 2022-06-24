package bot

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/config"
	"github.com/ggresillion/discordsoundboard/backend/internal/discord"
	"github.com/ggresillion/discordsoundboard/backend/internal/messaging"
)

type Bot struct {
	hub         *messaging.Hub
	session     *discordgo.Session
	ready       bool
	guildVoices map[string]*Audio
}

func NewBot(hub *messaging.Hub) *Bot {
	return &Bot{hub: hub, session: nil, guildVoices: make(map[string]*Audio), ready: false}
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
	b.session.AddHandler(b.commandHandler)

	b.session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuildMembers

	err = b.session.Open()
	if err != nil {
		return fmt.Errorf("error opening Discord session: %w", err)
	}

	b.ready = true

	log.Println("bot is now running")
	return nil
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	s.UpdateGameStatus(0, "!dsb")
}

func (b *Bot) guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	b.GetGuildVoice(event.Guild.ID).voiceStates = event.VoiceStates
}

func (b *Bot) GetGuildVoice(guildId string) *Audio {
	gv := b.guildVoices[guildId]
	if gv == nil {
		gv = NewAudio(guildId, b)
		b.guildVoices[guildId] = gv
	}
	return gv
}

func (b *Bot) GetLatency() time.Duration {
	return b.session.HeartbeatLatency()
}

func (b *Bot) GetGuilds(token string) ([]discord.Guild, error) {
	dc := discord.NewDiscordClient(token)
	userGuilds, err := dc.GetGuilds()
	if err != nil {
		return nil, err
	}

	guilds := make([]discord.Guild, 0)
	for _, g1 := range b.session.State.Guilds {
		for _, g2 := range userGuilds {
			if g1.ID == g2.ID {
				guilds = append(guilds, g2)
			}
		}
	}
	return guilds, nil
}
