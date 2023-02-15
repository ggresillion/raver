package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/display"
)

// Select this channel to hold the playlist
func (c *CommandHandler) playlist() *bot.CommandAndHandler {
	return &bot.CommandAndHandler{
		Command: &discordgo.ApplicationCommand{
			Name:        "playlist",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Display the current playlist in this channel",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// Send temporary interaction response
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Data: &discordgo.InteractionResponseData{
					Flags: discordgo.MessageFlagsEphemeral,
				},
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			if err != nil {
				handleError(s, i, err, "Failed to send temporary interaction response")
				return
			}

			// Get the music player
			m, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				handleError(s, i, err, "Failed to get music player")
				return
			}

			// Clear all messages in this channel
			// err = display.ClearMessages(s, i.ChannelID)
			// if err != nil {
			// 	handleError(s, i, err, "Failed to clear messages")
			// 	return
			// }

			// Print the playlist in this channel
			err = display.PrintPlaylist(s, i.GuildID, i.ChannelID, m.GetState())
			if err != nil {
				handleError(s, i, err, "Failed to print the playlist")
				return
			}

			// Subscribe to further changes to the playlist and print it
			subscriber := m.SubscribeToPlayerState()
			go func() {
				state := <-subscriber.Change
				display.PrintPlaylist(s, i.GuildID, i.ChannelID, state)
			}()

			// Clear interaction response
			err = s.InteractionResponseDelete(i.Interaction)
			if err != nil {
				handleError(s, i, err, "Failed to clear interaction response")
				return
			}
		},
	}
}
