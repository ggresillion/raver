package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal"
	"github.com/ggresillion/discordsoundboard/backend/internal/config"
)

type Bot struct {
	hub     *internal.Hub
	session *discordgo.Session
}

func NewBot(hub *internal.Hub) *Bot {
	return &Bot{hub: hub, session: nil}
}

var buffer = make([][]byte, 0)

func (b *Bot) StartBot() error {

	token := config.Get().BotToken

	if token == "" {
		return fmt.Errorf("no token provided")
	}

	var err error

	// Create a new Discord session using the provided bot token.
	b.session, err = discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	// Register ready as a callback for the ready events.
	b.session.AddHandler(ready)

	// Register messageCreate as a callback for the messageCreate events.
	b.session.AddHandler(messageCreate)

	// Register guildCreate as a callback for the guildCreate events.
	b.session.AddHandler(guildCreate)

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	b.session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open the websocket and begin listening.
	err = b.session.Open()
	if err != nil {
		return fmt.Errorf("error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("bot is now running")
	return nil
}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateGameStatus(0, "!dsb")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// check if the message is "!airhorn"
	if strings.HasPrefix(m.Content, "p") {
		handlePlaySoundMessage(s, m)
	}
}

// This function will be called (due to AddHandler above) every time a new
// guild is joined.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			// _, _ = s.ChannelMessageSend(channel.ID, "Airhorn is ready! Type !airhorn while in a voice channel to play a sound.")
			return
		}
	}
}

func handlePlaySoundMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Find the channel that the message came from.
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		return
	}

	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			// err = playSound(g.ID, vs.ChannelID)
			if err != nil {
				fmt.Println("Error playing sound:", err)
			}

			return
		}
	}
}

// playSound plays the current buffer to the provided channel.
// func playSound(guildID, channelID string) (err error) {

// 	// Join the provided voice channel.
// 	vc, err := Session.ChannelVoiceJoin(guildID, channelID, false, true)
// 	if err != nil {
// 		return err
// 	}

// 	youtube.PlayFromYoutube(vc, "tUlthCngK9U")

// 	return nil
// }

func (b *Bot) GetGuild(id string) (*discordgo.Guild, error) {
	return b.session.Guild(id)
}
