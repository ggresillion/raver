package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"raver/audio"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	token   string
	session *discordgo.Session
	gbots   map[string]*GBot
}

type GBot struct {
	Player            *audio.Player
	PlaylistMessageID string
	PlaylistChannelID string
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
	b.session.AddHandler(func(_ *discordgo.Session, _ *discordgo.Ready) { log.Println("bot: connected") })
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		g, err := b.Guild(i.GuildID)
		if err != nil {
			log.Println(err)
			return
		}
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			command := i.ApplicationCommandData().Name
			log.Printf("bot: received slash command: %s", command)
			if h, ok := CommandHandlers[command]; ok {
				h(g, s, i)
			}
		case discordgo.InteractionMessageComponent:
			command := i.MessageComponentData().CustomID
			log.Printf("bot: received component action: %s", command)
			if h, ok := CommandHandlers[command]; ok {
				h(g, s, i)
			}
		default:
			command := i.ApplicationCommandData().Name
			log.Printf("bot: received interaction: %s", command)
			if h, ok := CommandHandlers[command]; ok {
				h(g, s, i)
			}
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

	// Handle cleanup
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		for _, v := range b.gbots {
			b.session.ChannelMessageDelete(v.PlaylistChannelID, v.PlaylistMessageID)
		}
		os.Exit(0)
	}()

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
		Player:  audio.NewPlayer(guildID),
	}
	b.gbots[guildID] = g
	return g, nil
}

// func (g *GBot) BindStream(audioStream *audio.AudioStream) error {
// 	if g.vc == nil {
// 		return errors.New("no voice connection available")
// 	}
// 	err := g.vc.Speaking(true)
// 	if err != nil {
// 		return err
// 	}
//
// 	audioStream.Out = g.vc.OpusSend
//
// 	go func() {
// 		<-audioStream.OnStop()
// 		log.Println("bot: got stream end signal")
// 		err := g.vc.Speaking(false)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}()
//
// 	log.Printf("bot: playing stream %p", audioStream)
// 	return nil
// }

func (g *GBot) JoinUserChannel(userID string) error {
	log.Printf("bot: trying to join voice channel for user %s", userID)
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
			go func() {
				for {
					data, ok := <-g.Player.LineOut
					if !ok {
						return
					}
					g.vc.OpusSend <- data
				}
			}()
			vc.Speaking(true)
			log.Printf("bot: joinned voice channel %s", vc.ChannelID)
			return nil
		}
	}
	return fmt.Errorf("bot: user %s not in a voice channel", userID)
}
