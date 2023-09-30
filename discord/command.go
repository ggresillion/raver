package discord

import (
	"fmt"
	"log"
	"raver/youtube"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "play",
			Description: "Add a music to the playlist",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "query",
					Description:  "Search query",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	}

	commandHandlers = map[string]func(b *Bot, s *discordgo.Session, i *discordgo.InteractionCreate){
		"play": func(b *Bot, s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				data := i.ApplicationCommandData()
				trackID := data.Options[0].StringValue()

				stream, err := youtube.GetStreamFromYoutube(trackID)
				if err != nil {
					log.Println(err)
					return
				}
				i.User.
					b.PlayStream(stream, i.GuildID, i.ChannelID)

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Added %q to queue", trackID),
					},
				})
				if err != nil {
					log.Println(err)
				}
			case discordgo.InteractionApplicationCommandAutocomplete:
				data := i.ApplicationCommandData()
				q := data.Options[0].StringValue()
				if len(q) < 3 {
					return
				}
				tracks, err := youtube.Search(data.Options[0].StringValue())
				if err != nil {
					log.Println(err)
				}

				var choices []*discordgo.ApplicationCommandOptionChoice
				for _, track := range tracks {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  fmt.Sprintf("%q - %q", track.Title, track.Artist),
						Value: track.ID,
					})
				}

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionApplicationCommandAutocompleteResult,
					Data: &discordgo.InteractionResponseData{
						Choices: choices, // This is basically the whole purpose of autocomplete interaction - return custom options to the user.
					},
				})
				if err != nil {
					log.Println(err)
				}
			}
		},
	}
)
