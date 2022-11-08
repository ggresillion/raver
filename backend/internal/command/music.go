package command

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/music"
)

var (
	messagesToDeleteByInteractionID = make(map[string][]string)
)

// Extract infos from received commands
// Parameters are separated by '_'
func extractInfosFromPlayCommand(command string) (interactionID string, songId string) {
	splited := strings.Split(command, "_")
	return splited[0], splited[2]
}

// Handle the play response
func (c *CommandHandler) handlePlay(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		action := i.MessageComponentData().CustomID
		if strings.Contains(action, "_play_") {
			interactionID, _ := extractInfosFromPlayCommand(action)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Added song to playlist",
				},
			})

			messagesToDelete := messagesToDeleteByInteractionID[interactionID]
			if len(messagesToDelete) == 0 {
				return
			}
			err := s.ChannelMessagesBulkDelete(i.ChannelID, messagesToDelete)
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
			data := i.ApplicationCommandData()
			if len(data.Options) == 0 {
				return
			}
			q := data.Options[0].StringValue()

			result, err := c.musicManager.Search(q, 0)
			if err != nil {
				respond(s, i, fmt.Sprintf("failed to search %s", err.Error()))
				return
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Here are some songs i found",
				},
			})

			if err != nil {
				log.Error(err)
				return
			}

			messagesToDeleteByInteractionID[i.ID] = make([]string, 0)
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
					log.Error(err)
					return
				}
				messagesToDeleteByInteractionID[i.ID] = append(messagesToDeleteByInteractionID[i.ID], m.ID)
			}
		},
	}
}

// Play a song
// It first adds it to the playlist, then start the stream if no song is currently playing
func (c *CommandHandler) play() *bot.CommandAndHandler {
	return &bot.CommandAndHandler{
		Command: &discordgo.ApplicationCommand{
			Name:        "play",
			Description: "Play a song",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "query",
					Description: "name, id or url for the song",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		}, Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			data := i.ApplicationCommandData()
			if len(data.Options) == 0 {
				return
			}
			q := data.Options[0].StringValue()

			p, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
				return
			}

			if p.BotAudio().Status() == bot.NotConnected {
				p.BotAudio().JoinUserChannel(i.Member.User.ID)
			}

			res, err := c.musicManager.Search(q, 0)
			t := res.Tracks[0]

			err = p.AddToPlaylist(t.ID, music.TrackElement)
			if err != nil {
				respond(s, i, fmt.Sprintf("error while adding song to playlist %s", err.Error()))
				return
			}

			err = p.Play()
			if err != nil {
				respond(s, i, fmt.Sprintf("error while playing the song %s", err.Error()))
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       t.Title,
							Description: t.Album.Name,
							Image: &discordgo.MessageEmbedImage{
								URL: t.Thumbnail,
							},
						},
					},
				},
			})
		},
	}
}

// Pause the song
func (c *CommandHandler) pause() *bot.CommandAndHandler {
	return &bot.CommandAndHandler{
		Command: &discordgo.ApplicationCommand{
			Name:        "pause",
			Description: "Pause the current song",
		}, Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			p, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
				return
			}
			p.Pause()
			respond(s, i, "Paused current song")
		},
	}
}

// Stops the current song
func (c *CommandHandler) stop() *bot.CommandAndHandler {
	return &bot.CommandAndHandler{
		Command: &discordgo.ApplicationCommand{
			Name:        "stop",
			Description: "Stop the current song",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			p, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
				return
			}
			p.Stop()
			respond(s, i, "Stoped current song")
		},
	}
}

// Clear the playlist
func (c *CommandHandler) clear() *bot.CommandAndHandler {
	return &bot.CommandAndHandler{
		Command: &discordgo.ApplicationCommand{
			Name:        "clear",
			Description: "Clear the playlist",
		}, Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			p, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
				return
			}
			p.ClearPlaylist()
			respond(s, i, "Stoped current song")
		},
	}
}

// Skip the current song
// Automatically play the next song if the playlist is not empty
func (c *CommandHandler) skip() *bot.CommandAndHandler {
	return &bot.CommandAndHandler{
		Command: &discordgo.ApplicationCommand{
			Name:        "skip",
			Description: "Skip the current song",
		}, Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			p, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
				return
			}
			p.Skip()
			respond(s, i, "Skipped song")
		},
	}
}

// Return the current playlist
func (c *CommandHandler) playlist() *bot.CommandAndHandler {
	return &bot.CommandAndHandler{
		Command: &discordgo.ApplicationCommand{
			Name:        "playlist",
			Description: "Show the current playlist",
		}, Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			p, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
				return
			}
			playlist := p.Playlist()
			sPlaylist := ""
			for _, track := range playlist {
				sPlaylist = sPlaylist + track.Title
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please choose something",
					Flags:   1 << 6,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "📜",
									},
									Label: "Documentation",
									Style: discordgo.LinkButton,
									URL:   "https://discord.com/developers/docs/interactions/message-components#buttons",
								},
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "🔧",
									},
									Label: "Discord developers",
									Style: discordgo.LinkButton,
									URL:   "https://discord.gg/discord-developers",
								},
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "🦫",
									},
									Label: "Discord Gophers",
									Style: discordgo.LinkButton,
									URL:   "https://discord.gg/7RuRrVHyXF",
								},
							},
						},
					},
				},
			})
		},
	}
}
