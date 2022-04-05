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
func NewCommandHandler(b *bot.Bot, musicManager *music.MusicPlayerManager) *CommandHandler {
	h := &CommandHandler{b, musicManager}
	h.bot.RegisterCommands([]*bot.CommandAndHandler{
		h.join(),
		h.leave(),
		h.play(),
		h.pause(),
		h.stop(),
		h.clear(),
		h.skip(),
		h.playlist(),
	})
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
