package bot

import (
	"strings"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

type CommandAndHandler struct {
	Command *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
	Actions map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var (
	commands []*CommandAndHandler
)

// Handle all incoming commands that have been registered
func (b *Bot) commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		for _, c := range commands {
			for a, h := range c.Actions {
				if strings.HasPrefix(i.MessageComponentData().CustomID, a) {
					h(s, i)
				}
			}
		}
		return
	}
	for _, c := range commands {
		if c.Command.Name == i.ApplicationCommandData().Name {
			c.Handler(s, i)
		}
	}
}

// Register all commands
func (b *Bot) RegisterCommands(c []*CommandAndHandler) {
	commands = c

	// Register all commands to discord
	createdCommands, err := b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, "", lo.Map(commands, func(item *CommandAndHandler, index int) *discordgo.ApplicationCommand { return item.Command }))
	if err != nil {
		log.Fatalf("error trying to register commands %s", err)
	}

	for _, c := range createdCommands {
		log.Printf("registered bot command %s", c.Name)
	}
}

// Register an aditional handler
func (b *Bot) RegisterAdditionalHandler(f func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	b.session.AddHandler(f)
}
