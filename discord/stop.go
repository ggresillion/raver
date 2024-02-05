package discord

import (
	"github.com/bwmarrin/discordgo"
)

type StopCommand struct{}

func (c StopCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "stop",
		Description: "Stop and clear the playlist",
		Type:        discordgo.ChatApplicationCommand,
	}
}

func (c StopCommand) Handler(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	g.Player.Stop()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
	})

	if err != nil {
		sendError(s, i.Interaction, err)
		return
	}
}
