package discord

import (
	"github.com/bwmarrin/discordgo"
)

type PauseCommand struct{}

func (c PauseCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "pause",
		Description: "Pause the current track",
		Type:        discordgo.ChatApplicationCommand,
	}
}

func (c PauseCommand) Handler(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	g.Playlist.Pause()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Paused",
		},
	})
	if err != nil {
		sendError(s, i.Interaction, err)
		return
	}
}
