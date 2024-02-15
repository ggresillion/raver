package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Name() string
	Command() *discordgo.ApplicationCommand
	Handler(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate)
}

var (
	Commands = []Command{PlayCommand{}, PauseCommand{}, SkipCommand{}, PlaylistCommand{}, StopCommand{}}
)

func sendError(s *discordgo.Session, i *discordgo.Interaction, err error) error {
	log.Printf("discord: sending error to client: %v", err)
	_, err = s.FollowupMessageCreate(i, false, &discordgo.WebhookParams{
		Content: err.Error(),
		Flags:   discordgo.MessageFlagsEphemeral,
	})
	return err
}
