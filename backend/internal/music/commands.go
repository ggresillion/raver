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
	m.bot.RegisterCommand(&discordgo.ApplicationCommand{
		Name:        "pause",
		Description: "Pause the current song",
	}, m.pause)
	m.bot.RegisterCommand(&discordgo.ApplicationCommand{
		Name:        "stop",
		Description: "Stop the current song",
	}, m.stop)
	m.bot.RegisterCommand(&discordgo.ApplicationCommand{
		Name:        "clear",
		Description: "Clear the playlist",
	}, m.clear)
	m.bot.RegisterCommand(&discordgo.ApplicationCommand{
		Name:        "skip",
		Description: "Skip the current song",
	}, m.skip)
	m.bot.RegisterCommand(&discordgo.ApplicationCommand{
		Name:        "playlist",
		Description: "Show the current playlist",
	}, m.playlist)
}

func (m *MusicPlayerManager) play(s *discordgo.Session, i *discordgo.InteractionCreate) {

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

func (m *MusicPlayerManager) pause(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p, err := m.GetPlayer(i.GuildID)
	if err != nil {
		respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
		return
	}
	p.Pause()
	respond(s, i, "Paused current song")
}

func (m *MusicPlayerManager) stop(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, err := m.GetPlayer(i.GuildID)
	if err != nil {
		respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
		return
	}
	respond(s, i, "Stoped current song")
}

func (m *MusicPlayerManager) skip(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p, err := m.GetPlayer(i.GuildID)
	if err != nil {
		respond(s, i, fmt.Sprintf("cannot get player for guildID %s", i.GuildID))
		return
	}
	p.Skip()
	respond(s, i, "Stoped current song")
}

func (m *MusicPlayerManager) clear(s *discordgo.Session, i *discordgo.InteractionCreate) {
}

func (m *MusicPlayerManager) playlist(s *discordgo.Session, i *discordgo.InteractionCreate) {
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, m string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: m,
		},
	})
}
