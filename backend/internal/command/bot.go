package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
)

// Join the user current voice channel
func (c *CommandHandler) join() *bot.CommandAndHandler {
	return &bot.CommandAndHandler{
		Command: &discordgo.ApplicationCommand{
			Name:        "join",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Join the current channel",
		}, Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guildID := i.GuildID
			userID := i.Member.User.ID
			err := c.bot.GetGuildVoice(guildID).JoinUserChannel(userID)
			if err != nil {
				respond(s, i, "Please join a voice channel first")
				return
			}
			respond(s, i, "I just joined your voice channel")
		},
	}
}

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
