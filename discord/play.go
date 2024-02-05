package discord

import (
	"fmt"

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
		addToPlaylist(g, s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		autocompleteSearch(s, i)
	}
}

func autocompleteSearch(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	q := data.Options[0].StringValue()
	if len(q) < 3 {
		return
	}
	tracks, err := youtube.Search(data.Options[0].StringValue(), 0)
	if err != nil {
		sendError(s, i.Interaction, err)
		return
	}

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, track := range tracks {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  fmt.Sprintf("%q - %q", truncate(track.Title, 50), truncate(track.Artist, 50)),
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

func addToPlaylist(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	trackID := data.Options[0].StringValue()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		sendError(s, i.Interaction, err)
		return
	}

	err = g.JoinUserChannel(i.Member.User.ID)
	if err != nil {
		sendError(s, i.Interaction, err)
		return
	}

	track, err := youtube.GetPlayableTrackFromYoutube(trackID)
	if err != nil {
		sendError(s, i.Interaction, err)
		return
	}

	go func() {
		change := g.Player.Change
		for {
			<-change
			g.PrintPlaylist(s, i.Interaction)
		}
	}()

	err = g.Player.Add(track)
	if err != nil {
		sendError(s, i.Interaction, err)
		return
	}

	err = s.InteractionResponseDelete(i.Interaction)
	if err != nil {
		sendError(s, i.Interaction, err)
		return
	}
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max]
	}
	return s
}
