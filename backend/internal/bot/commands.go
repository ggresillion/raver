package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
)

func (b *Bot) commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}

func (b *Bot) RegisterCommand(command *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", command)
	if err != nil {
		log.Fatalf("error trying to register command %s %s", command.Name, err)
	}
	commandHandlers[command.Name] = handler
	log.Printf("registered bot command %s", command.Name)
}
