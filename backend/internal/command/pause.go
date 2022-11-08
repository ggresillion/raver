package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
)

// Pause the song
func (c *CommandHandler) pause() *bot.CommandAndHandler {
	return &bot.CommandAndHandler{
		Command: &discordgo.ApplicationCommand{
			Name:        "pause",
			Description: "Pause the current song",
		}, Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			p, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
				return
			}
			p.Pause()
			respond(s, i, "Paused current song")
		},
	}
}
