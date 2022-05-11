package bot

import (
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

type CommandAndHandler struct {
	Command *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var (
	commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
)

func (b *Bot) commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}

func (b *Bot) RegisterCommands(commandsAndHandlers []*CommandAndHandler) {
	var commands []*discordgo.ApplicationCommand

	for _, c := range commandsAndHandlers {
		commands = append(commands, c.Command)
		commandHandlers[c.Command.Name] = c.Handler
	}

	createdCommands, err := b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, "", commands)
	if err != nil {
		log.Fatalf("error trying to register commands %s", err)
	}

	for _, c := range createdCommands {
		log.Printf("registered bot command %s", c.Name)
	}
}
