package discord

import (
	"raver/audio"

	"github.com/bwmarrin/discordgo"
)

type PauseCommand struct{}

func (c PauseCommand) Name() string { return "pause" }

func (c PauseCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "pause",
		Description: "Pause the current track",
	}
}

func (c PauseCommand) Handler(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if g.Player.State == audio.Playing {
		g.Player.Pause()
	} else {
		g.Player.Resume()
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
	})
	if err != nil {
		sendError(s, i.Interaction, err)
		return
	}
}
