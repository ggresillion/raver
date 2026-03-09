package discord

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"raver/audio"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	token   string
	session *discordgo.Session
	gbots   map[string]*GBot
}

type GBot struct {
	Player                   *audio.Player
	PlaylistAlreadyDisplayed bool
	Guild                    *discordgo.Guild
	session                  *discordgo.Session
}

func NewBot(token string) *Bot {
	return &Bot{token: token, gbots: map[string]*GBot{}}
}

func (b *Bot) Session() *discordgo.Session {
	return b.session
}

func (b *Bot) Connect() error {
	if b.session == nil {
		// Create a new Discord session using the provided bot token.
		session, err := discordgo.New("Bot " + b.token)
		if err != nil {
			return fmt.Errorf("error creating Discord session: %w", err)
		}
		b.session = session
	}

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	b.session.Identify.Intents |= discordgo.IntentGuilds
	b.session.Identify.Intents |= discordgo.IntentGuildVoiceStates

	// Register handlers
	b.session.AddHandler(func(_ *discordgo.Session, _ *discordgo.Ready) { slog.Info("[bot] connected") })
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		g, err := b.Guild(i.GuildID)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		var command string
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			command = i.ApplicationCommandData().Name
			slog.Info("[bot] received slash command", "command", command, "guild_id", g.Guild.ID)
		case discordgo.InteractionMessageComponent:
			command = i.MessageComponentData().CustomID
			slog.Info("[bot] received component action", "command", command, "guild_id", g.Guild.ID)
		case discordgo.InteractionApplicationCommandAutocomplete:
			command = i.ApplicationCommandData().Name
			slog.Info("[bot] received autocomplete request", "command", command, "guild_id", g.Guild.ID)
		default:
			slog.Info("[bot] received unknown interaction", "interaction", i.Type.String(), "command", command, "guild_id", g.Guild.ID)
			return
		}
		handleCommand(command, g, s, i)
	})

	// Open the websocket and begin listening.
	err := b.session.Open()
	if err != nil {
		return fmt.Errorf("error opening Discord session: %w", err)
	}

	commands := make([]*discordgo.ApplicationCommand, 0)
	for _, c := range Commands {
		commands = append(commands, c.Command())
	}

	// Create commands
	_, err = b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, "", commands)
	if err != nil {
		return errors.Join(err, errors.New("cannot register commands"))
	}

	return nil
}

func (b *Bot) Stop() {
	for _, g := range b.gbots {
		g.Player.Stop()
	}
	b.session.Close()
}

func (b *Bot) Guild(guildID string) (*GBot, error) {
	g, exists := b.gbots[guildID]
	if exists {
		return g, nil
	}
	guild, err := b.session.State.Guild(guildID)
	if err != nil {
		return nil, fmt.Errorf("[bot] error getting guild %s: %v", guildID, err)
	}
	g = &GBot{
		session: b.session,
		Guild:   guild,
		Player:  audio.NewPlayer(guildID),
	}
	b.gbots[guildID] = g
	return g, nil
}

func (g *GBot) JoinUserChannel(userID string) error {
	slog.Info("[bot] trying to join voice channel for user", "user_id", userID)
	state, err := g.session.State.VoiceState(g.Guild.ID, userID)
	if err != nil {
		return fmt.Errorf("[bot] error getting voice state: %v", err)
	}

	vc, err := g.session.ChannelVoiceJoin(g.Guild.ID, state.ChannelID, false, true)
	if err != nil {
		return fmt.Errorf("[bot] error joining voice channel: %v", err)
	}

	g.Player.Plug(vc.OpusSend)
	vc.Speaking(true)
	slog.Info("[bot] joined voice channel", "channel_id", state.ChannelID, "guild_id", g.Guild.ID)
	return nil
}

func (g *GBot) LeaveChannel(ctx context.Context) error {
	vc, ok := g.session.VoiceConnections[g.Guild.ID]
	if !ok {
		return errors.New("no voice connection found")
	}
	vc.Disconnect()
	return nil
}

func handleCommand(command string, g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, v := range Commands {
		if v.Name() == command {
			v.Handler(g, s, i)
		}
	}
}
