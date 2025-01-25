package discord

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Name() string
	Command() *discordgo.ApplicationCommand
	Handler(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate)
}

var Commands = []Command{
	PlayCommand{},
	PauseCommand{},
	ResumeCommand{},
	SkipCommand{},
	PlaylistCommand{},
	StopCommand{},
}

func sendError(s *discordgo.Session, i *discordgo.Interaction, err error) error {
	slog.Error("discord: sending error to client", "error", err.Error(), "guild_id", i.GuildID)
	_, err = s.FollowupMessageCreate(i, false, &discordgo.WebhookParams{
		Content: err.Error(),
		Flags:   discordgo.MessageFlagsEphemeral,
	})
	return err
}
