package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	commands        []*discordgo.ApplicationCommand                                       = []*discordgo.ApplicationCommand{}
	commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
)

func (b *Bot) RegisterCommand(command *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	commands = append(commands, command)
	commandHandlers[command.Name] = handler

	_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", command)
	if err != nil {
		log.Fatalf("error trying to register command %s %s", command.Name, err)
	}
	log.Printf("registered bot command %s", command.Name)
}

func (b *Bot) registerCommands() {

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "join",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Join the current channel",
		},
	}

	commandHandlers["join"] = b.join

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

func (b *Bot) join(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildID := i.GuildID
	userID := i.Member.User.ID
	_, err := b.JoinUserChannel(guildID, userID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please join a voice channel first",
			},
		})
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "I just joined your voice channel",
		},
	})
}
