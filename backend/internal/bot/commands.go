package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	commands        []*discordgo.ApplicationCommand                                       = []*discordgo.ApplicationCommand{}
	commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
)

func (b *Bot) registerCommands() {
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	for _, c := range commands {
		_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", c)
		if err != nil {
			log.Fatalf("error trying to register command %s %s", c.Name, err)
		}
		log.Printf("registered bot command %s", c.Name)
	}
}

func (b *Bot) RegisterCommand(command *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	commands = append(commands, command)
	commandHandlers[command.Name] = handler
	if b.ready {
		_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", command)
		if err != nil {
			log.Fatalf("error trying to register command %s %s", command.Name, err)
		}
		log.Printf("registered bot command %s", command.Name)
	}
}
