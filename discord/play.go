package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"

	"raver/youtube"
)

type PlayCommand struct{}

func (c PlayCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "play",
		Description: "Add a music to the playlist",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "query",
				Description:  "Search query",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				Autocomplete: true,
			},
		},
	}
}

func (c PlayCommand) Handler(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()
		trackID := data.Options[0].StringValue()

		err := g.JoinUserChannel(i.Member.User.ID)
		if err != nil {
			sendError(s, i.Interaction, err)
			return
		}

		track, err := youtube.GetPlayableTrackFromYoutube(trackID)
		if err != nil {
			sendError(s, i.Interaction, err)
			return
		}

		g.Playlist.RegisterOnChange(func() { g.PrintPlaylist(i.Interaction) })

		err = g.Playlist.Add(track)
		if err != nil {
			sendError(s, i.Interaction, err)
			return
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Added %q to queue", trackID),
			},
		})
		if err != nil {
			sendError(s, i.Interaction, err)
			return
		}

	case discordgo.InteractionApplicationCommandAutocomplete:
		data := i.ApplicationCommandData()
		q := data.Options[0].StringValue()
		if len(q) < 3 {
			return
		}
		tracks, err := youtube.Search(data.Options[0].StringValue(), 0)
		if err != nil {
			log.Println(err)
			sendError(s, i.Interaction, err)
			return
		}

		var choices []*discordgo.ApplicationCommandOptionChoice
		for _, track := range tracks {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  fmt.Sprintf("%q - %q", track.Title, track.Artist),
				Value: track.ID,
			})
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
		if err != nil {
			sendError(s, i.Interaction, err)
			return
		}
	}

}
