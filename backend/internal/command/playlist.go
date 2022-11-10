package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/music"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

var (
	playlistMessages = make(map[string]string)
)

func clearMessages(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	messages, err := s.ChannelMessages(i.ChannelID, 100, "", "", "")
	if err != nil {
		return err
	}

	return s.ChannelMessagesBulkDelete(
		i.ChannelID,
		lo.Map(messages, func(m *discordgo.Message, i int) string { return m.ID }))
}

// Send a message containing the playlist
func printPlaylist(s *discordgo.Session, i *discordgo.InteractionCreate, state music.PlayerState) {

	embeds := lo.Map(state.Playlist, func(item *music.Track, index int) *discordgo.MessageEmbed {
		return &discordgo.MessageEmbed{
			Title:       item.Title,
			Description: item.Album.Name,
			Image: &discordgo.MessageEmbedImage{
				URL:    item.Thumbnail,
				Height: 100,
				Width:  100,
			},
		}
	})

	// Buttons
	var playPauseButton discordgo.MessageComponent

	switch state.Status {
	case bot.Playing:

		playPauseButton = discordgo.Button{
			Label: "",
			Emoji: discordgo.ComponentEmoji{
				Name: "⏸️",
			},
			Style:    discordgo.SecondaryButton,
			Disabled: false,
			CustomID: i.ID + "_pause",
		}

	case bot.Paused:
		playPauseButton = discordgo.Button{
			Label: "",
			Emoji: discordgo.ComponentEmoji{
				Name: "▶️",
			},
			Style:    discordgo.SecondaryButton,
			Disabled: false,
			CustomID: i.ID + "_play",
		}
	case bot.IDLE:
	default:
		playPauseButton = discordgo.Button{
			Label: "",
			Emoji: discordgo.ComponentEmoji{
				Name: "▶️",
			},
			Style:    discordgo.SecondaryButton,
			Disabled: true,
			CustomID: i.ID + "_play",
		}
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				playPauseButton,
				discordgo.Button{
					Label: "",
					Emoji: discordgo.ComponentEmoji{
						Name: "⏮️",
					},
					Style:    discordgo.SecondaryButton,
					Disabled: true,
					CustomID: i.ID + "_previous",
				},
				discordgo.Button{
					Label: "",
					Emoji: discordgo.ComponentEmoji{
						Name: "⏩",
					},
					Style:    discordgo.SecondaryButton,
					Disabled: true,
					CustomID: i.ID + "_skip",
				},
			},
		},
	}

	playlistID := playlistMessages[i.GuildID]
	if playlistID != "" {
		_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel:    i.ChannelID,
			ID:         playlistID,
			Embeds:     embeds,
			Components: components,
		})
		if err != nil {
			log.Error(err)
			return
		}
		return
	}

	m, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content:    "Playlist",
		Embeds:     embeds,
		Components: components,
	})
	if err != nil {
		log.Error(err)
		return
	}
	playlistMessages[i.GuildID] = m.ID
}

// Select this channel to hold the playlist
func (c *CommandHandler) playlist() *bot.CommandAndHandler {
	return &bot.CommandAndHandler{
		Command: &discordgo.ApplicationCommand{
			Name:        "playlist",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Display the current playlist in this channel",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			m, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				log.Error(err)
				return
			}
			clearMessages(s, i)
			printPlaylist(s, i, m.GetState())
			subscriber := m.SubscribeToPlayerState()
			go func() {
				state := <-subscriber.Change
				printPlaylist(s, i, state)
			}()
		},
	}
}
