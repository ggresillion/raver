package discord

import (
	"github.com/bwmarrin/discordgo"
)

type ResumeCommand struct{}

func (c ResumeCommand) Name() string { return "resume" }

func (c ResumeCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "resume",
		Description: "Resume the current track",
		Type:        discordgo.ChatApplicationCommand,
	}
}

func (c ResumeCommand) Handler(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	g.Player.Resume()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
	})
	if err != nil {
		sendError(s, i.Interaction, err)
		return
	}
}
