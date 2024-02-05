package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Command() *discordgo.ApplicationCommand
	Handler(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate)
}

var (
	Commands        = []*discordgo.ApplicationCommand{PlayCommand{}.Command()}
	CommandHandlers = map[string]func(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate){
		"play":     PlayCommand{}.Handler,
		"playlist": PlaylistCommand{}.Handler,
		"pause":    PauseCommand{}.Handler,
		"resume":   PauseCommand{}.Handler,
		"skip":     SkipCommand{}.Handler,
	}
)

func sendError(s *discordgo.Session, i *discordgo.Interaction, err error) error {
	log.Printf("discord: sending error to client: %v", err)
	_, err = s.FollowupMessageCreate(i, false, &discordgo.WebhookParams{
		Content: err.Error(),
		Flags:   discordgo.MessageFlagsEphemeral,
	})
	return err
	// return s.InteractionRespond(i, &discordgo.InteractionResponse{
	// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 	Data: &discordgo.InteractionResponseData{
	// 		Content: err.Error(),
	// 	},
	// })
}
