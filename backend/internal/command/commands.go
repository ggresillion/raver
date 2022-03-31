package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/music"
)

type CommandHandler struct {
	bot          *bot.Bot
	musicManager *music.MusicPlayerManager
}

// Command handler is responsible for handling bot slash commands
func NewCommandHandler(bot *bot.Bot, musicManager *music.MusicPlayerManager) *CommandHandler {
	h := &CommandHandler{bot, musicManager}
	h.bot.RegisterCommand(h.join())
	h.bot.RegisterCommand(h.leave())
	h.bot.RegisterCommand(h.play())
	h.bot.RegisterCommand(h.pause())
	h.bot.RegisterCommand(h.stop())
	h.bot.RegisterCommand(h.clear())
	h.bot.RegisterCommand(h.skip())
	h.bot.RegisterCommand(h.playlist())
	return h
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, m string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: m,
		},
	})
}
