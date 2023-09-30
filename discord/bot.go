package discord

import (
	"fmt"
	"log"
	"raver/audio"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	token   string
	session *discordgo.Session
}

func NewBot(token string) *Bot {
	return &Bot{token: token}
}

func (b *Bot) Connect() error {
	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + b.token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}
	b.session = session

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	b.session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Register handlers
	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is connected") })
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(b, s, i)
		}
	})

	// Open the websocket and begin listening.
	err = b.session.Open()
	if err != nil {
		return fmt.Errorf("error opening Discord session: %w", err)
	}

	// Create commands
	_, err = b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, "", commands)
	if err != nil {
		log.Fatalf("Cannot register commands: %v", err)
	}

	return nil
}

func (b *Bot) Stop() {
	b.session.Close()
}

func (b *Bot) PlayStream(audioStream *audio.AudioStream, guildID, channelID string) (err error) {

	log.Printf("Connecting to channel %q in guild %q", channelID, guildID)

	vc, err := b.session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	vc.Speaking(true)

	go func() {
		audioStream.Out = vc.OpusSend
		audioStream.Play()
		audioStream.BlockUntilStop()
		vc.Speaking(false)
	}()

	return nil
}
