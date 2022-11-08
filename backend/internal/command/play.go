package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/music"
	log "github.com/sirupsen/logrus"
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
			if err != nil {
				log.Error(err)
				return
			}

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
