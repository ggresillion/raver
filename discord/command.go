package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Command() *discordgo.ApplicationCommand
	Handler(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate)
}

var (
	Commands        = []*discordgo.ApplicationCommand{PlayCommand{}.Command()}
	CommandHandlers = map[string]func(g *GBot, s *discordgo.Session, i *discordgo.InteractionCreate){
		"play": PlayCommand{}.Handler,
	}
)

func sendError(s *discordgo.Session, i *discordgo.Interaction, err error) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: err.Error(),
		},
	})
}
