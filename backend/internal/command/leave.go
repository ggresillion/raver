package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
)

// Leave the voice channel
func (c *CommandHandler) leave() *bot.CommandAndHandler {
	return &bot.CommandAndHandler{
		Command: &discordgo.ApplicationCommand{
			Name:        "leave",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Leave the channel",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guildID := i.GuildID
			c.bot.GetGuildVoice(guildID).LeaveChannel()
			respond(s, i, "I left the channel")
		},
	}
}
