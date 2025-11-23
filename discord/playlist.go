package discord

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"raver/audio"

	"github.com/bwmarrin/discordgo"
)

/**
Button codes
play 1142048243073695765
pause 892404982970740786
skip 892404982911991949
previous 892404982538719274
forward 924991555771711518
backwards 924991555847200778
volume down 924991555893342268
volume up 924991555842998333
replay 892404982521933876
replay playlist 1016805948779663501
stop 892404982828113961
*/

const sizeProgressBar = 27

type PlaylistCommand struct{}

func (c PlaylistCommand) Name() string { return "playlist" }

func (c PlaylistCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "playlist",
		Description: "Display the playlist in this channel",
	}
}

func (c PlaylistCommand) Handler(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	g.PrintPlaylist(s, i.Interaction)
}

func (g *GBot) PrintPlaylist(s *discordgo.Session, i *discordgo.Interaction) {
	if g.PlaylistAlreadyDisplayed {
		return
	}
	slog.Info("bot: creating playlist display", "guild_id", i.GuildID)
	m, err := g.session.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content:    "",
		Embeds:     generateEmbeds(g),
		Components: generateComponents(g),
	})
	if err != nil {
		sendError(s, i, err)
		return
	}
	g.PlaylistAlreadyDisplayed = true
	go func() {
		defer func() {
			slog.Info("bot: deleting playlist display", "guild_id", i.GuildID)
			err := s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				sendError(s, i, err)
				return
			}
			g.PlaylistAlreadyDisplayed = false
			err = g.vc.Disconnect(context.Background())
			if err != nil {
				sendError(s, i, err)
				return
			}
			return
		}()
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-g.Player.Change:
				slog.Info("bot: got playlist update", "guild_id", i.GuildID, "tracks", len(g.Player.Queue), "state", g.Player.State)
				if len(g.Player.Queue) == 0 {
					return
				}
			case <-ticker.C:
				if len(g.Player.Queue) == 0 {
					slog.Info("bot: playlist empty, stoping update loop", "guild_id", i.GuildID)
					return
				}
				embeds := generateEmbeds(g)
				components := generateComponents(g)
				_, err := g.session.ChannelMessageEditComplex(&discordgo.MessageEdit{
					ID:         m.ID,
					Channel:    m.ChannelID,
					Embeds:     &embeds,
					Components: &components,
				})
				if err != nil {
					sendError(s, i, err)
					return
				}
			}
		}
	}()
}

func generateEmbeds(g *GBot) []*discordgo.MessageEmbed {
	if len(g.Player.Queue) == 0 {
		return []*discordgo.MessageEmbed{}
	}
	var fields []*discordgo.MessageEmbedField
	for _, track := range g.Player.Queue {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%s [%s]", track.Title, formatDuration(track.Duration)),
			Value: track.Artist,
		})
	}
	var progressBar []rune
	for i := 0; i < sizeProgressBar; i++ {
		progressBar = append(progressBar, '▬')
	}
	progressIndex := int(float64(g.Player.Progress()) / 100 * sizeProgressBar)
	progressBar[progressIndex] = '🔘'

	fields = append(fields, &discordgo.MessageEmbedField{
		Name: string(progressBar),
	})
	return []*discordgo.MessageEmbed{
		{
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", g.Player.Queue[0].ID),
			},
			Title:  "Now playing",
			Type:   discordgo.EmbedTypeRich,
			Fields: fields,
		},
	}
}

func generateComponents(g *GBot) []discordgo.MessageComponent {
	var playPause discordgo.Button
	if g.Player.State == audio.Playing {
		playPause = discordgo.Button{
			CustomID: "pause",
			Style:    discordgo.SecondaryButton,
			Emoji: &discordgo.ComponentEmoji{
				Name: "custom",
				ID:   "1142048234123042907",
			},
		}
	} else {
		playPause = discordgo.Button{
			CustomID: "resume",
			Style:    discordgo.SecondaryButton,
			Emoji: &discordgo.ComponentEmoji{
				Name: "custom",
				ID:   "1142048243073695765",
			},
		}
	}

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				playPause,
				discordgo.Button{
					CustomID: "skip",
					Style:    discordgo.SecondaryButton,
					Disabled: len(g.Player.Queue) < 2,
					Emoji: &discordgo.ComponentEmoji{
						Name: "custom",
						ID:   "1142048259607629905",
					},
				},
				discordgo.Button{
					CustomID: "stop",
					Style:    discordgo.DangerButton,
					Emoji: &discordgo.ComponentEmoji{
						Name: "custom",
						ID:   "892404982828113961",
					},
				},
			},
		},
	}
}
