package command

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/music"
	log "github.com/sirupsen/logrus"
)

var (
	messagesToDeleteByInteractionID = make(map[string][]string)
)

// Extract infos from received commands
// Parameters are separated by '_'
func extractInfosFromPlayCommand(command string) (interactionID string, trackID string) {
	splited := strings.Split(command, "_")
	return splited[0], splited[2]
}

// Handle the play response
func (c *CommandHandler) handlePlayResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		action := i.MessageComponentData().CustomID
		if strings.Contains(action, "_play_") {
			interactionID, trackID := extractInfosFromPlayCommand(action)

			// Add to playlist
			player, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				log.Error(err)
				return
			}
			err = player.AddToPlaylist(trackID, music.TrackElement)
			if err != nil {
				log.Error(err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}

			// Delete track messages
			messagesToDelete := messagesToDeleteByInteractionID[interactionID]
			if len(messagesToDelete) == 0 {
				return
			}
			err = s.ChannelMessagesBulkDelete(i.ChannelID, messagesToDelete)
			if err != nil {
				log.Error(err)
				return
			}
		}
	}
}

// Search for a song
func (c *CommandHandler) search() *bot.CommandAndHandler {
	return &bot.CommandAndHandler{
		Command: &discordgo.ApplicationCommand{
			Name:        "search",
			Description: "Search for a song",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "query",
					Description: "name, id or url for the song",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		}, Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Get params
			data := i.ApplicationCommandData()
			if len(data.Options) == 0 {
				return
			}
			q := data.Options[0].StringValue()

			// Search for the songs
			result, err := c.musicManager.Search(q, 0)
			if err != nil {
				handleError(s, i, err, "failed to search")
				return
			}

			// Respond to the interaction
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "Here are the songs I found: ",
				},
			})
			if err != nil {
				handleError(s, i, err, "failed to respond")
				return
			}

			err = s.InteractionResponseDelete(i.Interaction)

			// Append tracks to response
			messagesToDeleteByInteractionID[i.ID] = []string{}
			for index, r := range result.Tracks {

				if index > 4 {
					return
				}

				artists := ""
				for _, a := range r.Artists {
					artists = artists + a.Name
				}

				m, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
					Content: r.Title,
					Embeds: []*discordgo.MessageEmbed{
						{
							Type:        discordgo.EmbedTypeImage,
							Title:       r.Title,
							Description: artists,
							URL:         "https://www.youtube.com/watch?v=" + r.ID,
							Image: &discordgo.MessageEmbedImage{
								URL: r.Thumbnail,
							},
						},
					},
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Label:    "Play",
									Style:    discordgo.SecondaryButton,
									Disabled: false,
									CustomID: i.ID + "_play_" + r.ID,
								},
							},
						},
					},
				})
				if err != nil {
					handleError(s, i, err, "failed to print search results")
					return
				}
				messagesToDeleteByInteractionID[i.ID] = append(messagesToDeleteByInteractionID[i.ID], m.ID)
			}
		},
	}
}
