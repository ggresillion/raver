package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/music"
)

type CommandHandler struct {
	bot          *bot.Bot
	musicManager *music.PlayerManager
}

// Command handler is responsible for handling bot slash commands
func NewCommandHandler(b *bot.Bot, musicManager *music.PlayerManager) *CommandHandler {
	c := &CommandHandler{b, musicManager}
	c.bot.RegisterCommands([]*bot.CommandAndHandler{
		c.join(),
		c.leave(),
		c.play(),
		c.pause(),
		c.stop(),
		c.clear(),
		c.skip(),
		c.playlist(),
		c.search(),
	})
	c.bot.RegisterAdditionalHandler(c.handlePlayResponse)
	return c
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, m string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: m,
		},
	})
}
