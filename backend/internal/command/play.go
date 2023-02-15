package command

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/music"
	"github.com/samber/lo"
)

// Play a song
// It first adds it to the playlist, then start the stream if no song is currently playing
func (c *CommandHandler) play() *bot.CommandAndHandler {
	return &bot.CommandAndHandler{
		Command: &discordgo.ApplicationCommand{
			Name:        "play",
			Description: "Play a song",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "query",
					Description:  "name, id or url for the song",
					Required:     true,
					Type:         discordgo.ApplicationCommandOptionString,
					Autocomplete: true,
				},
			},
		}, Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				data := i.ApplicationCommandData()
				ID := data.Options[0].StringValue()
				p, err := c.musicManager.GetPlayer(i.GuildID)
				if err != nil {
					respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
					return
				}

				if p.BotAudio().Status() == bot.NotConnected {
					err = p.BotAudio().JoinUserChannel(i.Member.User.ID)
					if err != nil {
						respond(s, i, "failed to join user channel")
						return
					}
				}

				track, err := p.AddSongToPlaylist(ID)
				if err != nil {
					respond(s, i, fmt.Sprintf("error while adding the song %s", err.Error()))
					return
				}
				p.Play()
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: fmt.Sprintf("Added %s - %s to playlist!", track.Title, formatTrackArtists(track)),
					},
				})
			case discordgo.InteractionApplicationCommandAutocomplete:
				data := i.ApplicationCommandData()
				query := data.Options[0].StringValue()
				if query == "" {
					return
				}
				res, err := c.musicManager.SearchSimple(query, 0)
				if err != nil {
					respond(s, i, fmt.Sprintf("error while fetching songs %s", err.Error()))
					return
				}

				choices := lo.Map(res, func(t *music.Track, i int) *discordgo.ApplicationCommandOptionChoice {
					return &discordgo.ApplicationCommandOptionChoice{
						Name:  fmt.Sprintf("%s - %s", t.Title, formatTrackArtists(t)),
						Value: t.ID,
					}
				})
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionApplicationCommandAutocompleteResult,
					Data: &discordgo.InteractionResponseData{
						Choices: choices,
					},
				})
				if err != nil {
					respond(s, i, fmt.Sprintf("error while getting song autocompletion %s", err.Error()))
					return
				}
			}
		},
	}
}

func formatTrackArtists(track *music.Track) string {
	return strings.Join(lo.Map(track.Artists, func(a music.Artist, i int) string { return a.Name }), ", ")
}
