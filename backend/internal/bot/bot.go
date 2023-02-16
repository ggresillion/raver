package bot

import (
	"fmt"
	"time"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	botToken    string
	session     *discordgo.Session
	ready       bool
	guildVoices map[string]*Audio
}

func NewBot(botToken string) *Bot {
	return &Bot{botToken: botToken, session: nil, guildVoices: make(map[string]*Audio), ready: false}
}

func (b *Bot) StartBot() error {

	token := b.botToken

	if token == "" {
		return fmt.Errorf("no token provided")
	}

	var err error
	b.session, err = discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	b.session.Debug = true

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

func (b *Bot) GetGuilds(token string) ([]*discordgo.UserGuild, error) {
	userGuilds, err := b.session.UserGuilds(100, "", "", discordgo.WithHeader("Authorization", "Bearer "+token))
	if err != nil {
		return nil, err
	}

	return lo.Filter(userGuilds, func(g1 *discordgo.UserGuild, _ int) bool {
		return lo.ContainsBy(b.session.State.Guilds, func(g2 *discordgo.Guild) bool {
			return g1.ID == g2.ID
		})
	}), nil
}

func (b *Bot) GetUser(token string) (*discordgo.User, error) {
	return b.session.User("@me", discordgo.WithHeader("Authorization", "Bearer "+token))
}
