package discord

import (
	"fmt"
	"log"
	"raver/audio"

	"github.com/bwmarrin/discordgo"
)

type PlaylistCommand struct{}

func (c PlaylistCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "playlist",
		Description: "Display the playlist in this channel",
		Type:        discordgo.ChatApplicationCommand,
	}
}

func (c PlaylistCommand) Handler(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	g.PrintPlaylist(s, i.Interaction)
}

func (g *GBot) PrintPlaylist(s *discordgo.Session, i *discordgo.Interaction) {
	log.Printf("bot: updating playlist display for: %s", i.ChannelID)
	var tracks []*discordgo.MessageEmbedField
	for _, track := range g.Player.Queue {
		tracks = append(tracks, &discordgo.MessageEmbedField{
			Name:  track.Title,
			Value: track.Artist,
		})
	}

	embeds := []*discordgo.MessageEmbed{
		{
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", g.Player.Queue[0].ID),
			},
			Title:  "Now playing",
			Type:   discordgo.EmbedTypeRich,
			Fields: tracks,
		},
	}

	var playPause discordgo.MessageComponent
	if g.Player.State == audio.Playing {
		playPause = discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: "pause",
					Style:    discordgo.SecondaryButton,
					Emoji: discordgo.ComponentEmoji{
						Name: "custom",
						ID:   "1142048234123042907",
					},
				},
			},
		}
	} else {
		playPause = discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: "resume",
					Style:    discordgo.SecondaryButton,
					Emoji: discordgo.ComponentEmoji{
						Name: "custom",
						ID:   "1142048243073695765",
					},
				},
			},
		}
	}

	components := []discordgo.MessageComponent{
		playPause,
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: "skip",
					Style:    discordgo.SecondaryButton,
					Disabled: len(g.Player.Queue) < 2,
					Emoji: discordgo.ComponentEmoji{
						Name: "custom",
						ID:   "1142048259607629905",
					},
				},
			},
		},
	}

	if g.PlaylistMessageID != "" {
		_, err := g.session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:         g.PlaylistMessageID,
			Channel:    i.ChannelID,
			Embeds:     embeds,
			Components: components,
		})
		// workaround because message was sometimes not updated
		_, err = g.session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:         g.PlaylistMessageID,
			Channel:    i.ChannelID,
			Embeds:     embeds,
			Components: components,
		})
		if err != nil {
			sendError(s, i, err)
			return
		}
	} else {
		m, err := g.session.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Content:    "",
			Embeds:     embeds,
			Components: components,
		})
		if err != nil {
			sendError(s, i, err)
			return
		}
		g.PlaylistMessageID = m.ID
		g.PlaylistChannelID = m.ChannelID
	}
}
