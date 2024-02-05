package discord

import (
	"github.com/bwmarrin/discordgo"
)

type SkipCommand struct{}

func (c SkipCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "skip",
		Description: "Skip the current track",
		Type:        discordgo.ChatApplicationCommand,
	}
}

func (c SkipCommand) Handler(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	g.Player.Skip()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
	})

	if err != nil {
		sendError(s, i.Interaction, err)
		return
	}
}
