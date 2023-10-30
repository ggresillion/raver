package discord

import (
	"errors"
	"fmt"
	"log"

	"raver/audio"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	token   string
	session *discordgo.Session
	gbots   map[string]*GBot
}

type GBot struct {
	Playlist          *audio.Playlist
	PlaylistMessageID *string
	session           *discordgo.Session
	guild             *discordgo.Guild
	vc                *discordgo.VoiceConnection
}

func NewBot(token string) *Bot {
	return &Bot{token: token, gbots: map[string]*GBot{}}
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
	b.session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Register handlers
	b.session.AddHandler(func(_ *discordgo.Session, _ *discordgo.Ready) { log.Println("Bot: connected") })
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			g, err := b.Guild(i.GuildID)
			if err != nil {
				log.Println(err)
				return
			}
			h(g, s, i)
		}
	})

	// Open the websocket and begin listening.
	err := b.session.Open()
	if err != nil {
		return fmt.Errorf("error opening Discord session: %w", err)
	}

	// Create commands
	_, err = b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, "", Commands)
	if err != nil {
		log.Fatalf("Cannot register commands: %v", err)
	}

	return nil
}

func (b *Bot) Stop() {
	b.session.Close()
}

func (b *Bot) Guild(guildID string) (*GBot, error) {
	g, exists := b.gbots[guildID]
	if exists {
		return g, nil
	}
	guild, err := b.session.State.Guild(guildID)
	if err != nil {
		return nil, err
	}
	g = &GBot{
		session: b.session,
		guild:   guild,
		Playlist: audio.NewPlaylist(
			guildID,
			func(stream *audio.AudioStream) error {
				return g.BindStream(stream)
			}),
	}
	b.gbots[guildID] = g
	return g, nil
}

func (g *GBot) BindStream(audioStream *audio.AudioStream) error {
	if g.vc == nil {
		return errors.New("no voice connection available")
	}
	err := g.vc.Speaking(true)
	if err != nil {
		return err
	}

	audioStream.Out = g.vc.OpusSend

	go func() {
		audioStream.BlockUntilStop()
		err := g.vc.Speaking(false)
		if err != nil {
			log.Println(err)
		}
	}()

	log.Printf("Bot: playing stream %p", audioStream)
	return nil
}

func (g *GBot) JoinUserChannel(userID string) error {
	log.Printf("Bot: trying to join voice channel for user %s", userID)
	// guild, err := g.session.State.Guild(g.guild.ID)
	// if err != nil {
	// 	return err
	// }
	for _, v := range g.guild.VoiceStates {
		if v.UserID == userID {
			vc, err := g.session.ChannelVoiceJoin(g.guild.ID, v.ChannelID, false, true)
			if err != nil {
				return err
			}

			g.vc = vc
			vc.Speaking(true)
			log.Printf("Bot: joinned voice channel %s", vc.ChannelID)
			return nil
		}
	}
	return fmt.Errorf("Bot: user %s not in a voice channel", userID)
}
