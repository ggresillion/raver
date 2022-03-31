package command

import (
	"github.com/bwmarrin/discordgo"
)

// Join the user current voice channel
func (c *CommandHandler) join() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:        "join",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Join the current channel",
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guildID := i.GuildID
			userID := i.Member.User.ID
			_, err := c.bot.JoinUserChannel(guildID, userID)
			if err != nil {
				respond(s, i, "Please join a voice channel first")
				return
			}
			respond(s, i, "I just joined your voice channel")
		}
}

// Leave the voice channel
func (c *CommandHandler) leave() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:        "leave",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Leave the channel",
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guildID := i.GuildID
			c.bot.LeaveChannel(guildID)
			respond(s, i, "I left the channel")
		}
}
