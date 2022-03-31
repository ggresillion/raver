package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Play a song
// It first adds it to the playlist, then start the stream if no song is currently playing
func (c *CommandHandler) play() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:        "play",
			Description: "Play a song",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "link",
					Description: "url or id for the song",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			data := i.ApplicationCommandData()
			if len(data.Options) == 0 {
				return
			}
			ID := data.Options[0].StringValue()

			p, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
				return
			}

			t, err := p.AddToPlaylist(ID)
			if err != nil {
				respond(s, i, fmt.Sprintf("error while adding song to playlist %s", err.Error()))
				return
			}

			err = p.Play()
			if err != nil {
				respond(s, i, fmt.Sprintf("error while playing the song %s", err.Error()))
				return
			}
			respond(s, i, fmt.Sprintf("Adding %s to queue", t.Title))
		}
}

// Pause the song
func (c *CommandHandler) pause() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:        "pause",
			Description: "Pause the current song",
		}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			p, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
				return
			}
			p.Pause()
			respond(s, i, "Paused current song")
		}
}

// Stops the current song
func (c *CommandHandler) stop() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:        "stop",
			Description: "Stop the current song",
		}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			_, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
				return
			}
			respond(s, i, "Stoped current song")
		}
}

// Clear the playlist
func (c *CommandHandler) clear() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
			Name:        "clear",
			Description: "Clear the playlist",
		}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			p, err := c.musicManager.GetPlayer(i.GuildID)
			if err != nil {
				respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
				return
			}
			p.Skip()
			respond(s, i, "Stoped current song")
		}
}

// Skip the current song
// Automatically play the next song if the playlist is not empty
func (c *CommandHandler) skip() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
		Name:        "skip",
		Description: "Skip the current song",
	}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {}
}

// Return the current playlist
func (c *CommandHandler) playlist() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	return &discordgo.ApplicationCommand{
		Name:        "playlist",
		Description: "Show the current playlist",
	}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {}
}
