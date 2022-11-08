package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
)

// Skip the current song
// Automatically play the next song if the playlist is not empty
func (c *CommandHandler) skip() *bot.CommandAndHandler {
	return &bot.CommandAndHandler{
		Command: &discordgo.ApplicationCommand{
			Name:        "skip",
			Description: "Skip the current song",
		}, Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			p, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
				return
			}
			p.Skip()
			respond(s, i, "Skipped song")
		},
	}
}
