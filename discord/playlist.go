package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func (g *GBot) PrintPlaylist(i *discordgo.Interaction) {

	var tracks []*discordgo.MessageEmbedField
	for _, track := range g.Playlist.Queue {
		tracks = append(tracks, &discordgo.MessageEmbedField{
			Name:  track.Title,
			Value: track.Artist,
		})
	}

	embeds := []*discordgo.MessageEmbed{
		{
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://www.startpage.com/av/proxy-image?piurl=http%3A%2F%2Fdroplr.com%2Fwp-content%2Fuploads%2F2020%2F10%2FDiscord-music-e1635364775454.png&sp=1696163152Tb7a0d1728f87144df9c9d0fed65754d282d2be0b09c87afdabaa4e29ba3cca27",
			},
			Title:  "Playlist",
			Type:   discordgo.EmbedTypeRich,
			Fields: tracks,
		},
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: "pause",
					Style:    discordgo.SecondaryButton,
					Emoji: discordgo.ComponentEmoji{
						Name: "custom",
						ID:   "892404982970740786",
					},
				},
			},
		},
	}

	if g.PlaylistMessageID != nil {
		_, err := g.session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:         *g.PlaylistMessageID,
			Channel:    i.ChannelID,
			Embeds:     embeds,
			Components: components,
		})
		if err != nil {
			log.Println(fmt.Errorf("could not print playlist: %w", err))
		}
	} else {
		m, err := g.session.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Content:    "",
			Embeds:     embeds,
			Components: components,
		})
		if err != nil {
			log.Println(fmt.Errorf("could not print playlist: %w", err))
		}
		g.PlaylistMessageID = &m.ID
	}

}
