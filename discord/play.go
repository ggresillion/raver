package discord

import (
	"fmt"
	"raver/youtube"
	"raver/youtube/ytdlp"
	"time"

	"github.com/bwmarrin/discordgo"
)

type PlayCommand struct{}

func (c PlayCommand) Name() string { return "play" }

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
	tracks, err := youtube.Search(i.GuildID, data.Options[0].StringValue(), 0)
	if err != nil {
		sendError(s, i.Interaction, err)
		return
	}

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, track := range tracks {
		name := fmt.Sprintf("%q - **%q** [%s]", truncate(track.Title, 50), truncate(track.Artist, 50), formatDuration(track.Duration))
		if len(name) < 1 || len(name) > 100 {
			continue
		}
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
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

	track, err := youtube.NewYoutube(ytdlp.NewYoutubeAdapter()).GetPlayableTrackFromYoutube(i.GuildID, trackID)
	if err != nil {
		sendError(s, i.Interaction, err)
		return
	}

	err = g.Player.Add(track)
	if err != nil {
		sendError(s, i.Interaction, err)
		return
	}

	g.PrintPlaylist(s, i.Interaction)

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

func formatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
