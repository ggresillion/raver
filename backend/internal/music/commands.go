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
	}, m.play)
}

func (m *MusicPlayerManager) play(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	if len(data.Options) == 0 {
		return
	}
	ID := data.Options[0].StringValue()

	p, err := m.GetPlayer(i.GuildID)
	if err != nil {
		respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
		return
	}

	t, err := p.AddToPlaylist(ID)
	if err != nil {
		respond(s, i, fmt.Sprintf("error while adding song to playlist %s", err.Error()))
		return
	}

	_, err = p.Play()
	if err != nil {
		respond(s, i, fmt.Sprintf("error while playing the song %s", err.Error()))
		return
	}
	respond(s, i, fmt.Sprintf("Adding %s to queue", t.Title))
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, m string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: m,
		},
	})
}
