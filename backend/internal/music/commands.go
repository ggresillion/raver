package music

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (m *MusicPlayerManager) registerCommands() {
	m.bot.RegisterCommand(&discordgo.ApplicationCommand{
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

		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		data := i.ApplicationCommandData()
		if len(data.Options) == 0 {
			return
		}
		ID := data.Options[0].StringValue()

		t, err := m.play(i.GuildID, ID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Adding %s to queue", t.Title),
			},
		})
	})
}

func (m *MusicPlayerManager) play(guildID, ID string) (*Track, error) {
	p, err := m.GetPlayer(guildID)
	if err != nil {
		return nil, fmt.Errorf("cannot get player for guildID %s", guildID)
	}

	t, err := p.AddToPlaylist(ID)
	if err != nil {
		return nil, fmt.Errorf("error while adding song to playlist %s", err)
	}

	err = p.Play()
	if err != nil {
		return nil, fmt.Errorf("error while playing the song %s", err)
	}
	return t, nil
}
