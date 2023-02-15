package display

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/music"
	"github.com/samber/lo"
)

var (
	playlistMessages = make(map[string]string)
)

func PrintPlaylist(s *discordgo.Session, guildID, channelID string, state music.PlayerState) error {

	// Buttons
	var playPauseButton discordgo.MessageComponent
	switch state.Status {
	case bot.Playing:
		playPauseButton = pauseButton(guildID)
	case bot.Paused:
		playPauseButton = playButton(guildID, false)
	case bot.IDLE:
	default:
		playPauseButton = playButton(guildID, true)
	}
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				playPauseButton,
				previousButton(guildID),
				skipButton(guildID),
			},
		},
	}

	// If the playlist already exists in this channel, updates it, else, creates it
	// playlistID := playlistMessages[guildID]
	// if playlistID != "" {
	// 	_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
	// 		Channel:    channelID,
	// 		ID:         playlistID,
	// 		Embeds:     tracks(state),
	// 		Components: components,
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }

	m, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: "",
		Embeds: []*discordgo.MessageEmbed{{
			Title: "Queue: " + strconv.Itoa(len(state.Playlist)),
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Playlist",
				IconURL: s.State.User.AvatarURL(""),
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: thumbnail(state),
			},
			Fields: tracks(state),
		},
		},
		Components: components,
	})
	if err != nil {
		return err
	}
	playlistMessages[guildID] = m.ID
	return nil
}

func playButton(ID string, disabled bool) *discordgo.Button {
	return &discordgo.Button{
		Label: "",
		Emoji: discordgo.ComponentEmoji{
			ID:   "892404982970740786",
			Name: "Play",
		},
		Style:    discordgo.PrimaryButton,
		Disabled: disabled,
		CustomID: ID + "_play",
	}
}

func pauseButton(ID string) *discordgo.Button {
	return &discordgo.Button{
		Label: "",
		Emoji: discordgo.ComponentEmoji{
			ID:   "892404982970740786",
			Name: "Pause",
		},
		Style:    discordgo.PrimaryButton,
		Disabled: false,
		CustomID: ID + "_pause",
	}
}

func previousButton(ID string) *discordgo.Button {
	return &discordgo.Button{
		Label: "",
		Emoji: discordgo.ComponentEmoji{
			ID:   "892404982538719274",
			Name: "Back",
		},
		Style:    discordgo.PrimaryButton,
		Disabled: true,
		CustomID: ID + "_previous",
	}
}

func skipButton(ID string) *discordgo.Button {
	return &discordgo.Button{
		Label: "",
		Emoji: discordgo.ComponentEmoji{
			ID:   "892404982911991949",
			Name: "Skip",
		},
		Style:    discordgo.PrimaryButton,
		Disabled: true,
		CustomID: ID + "_skip",
	}
}

func tracks(state music.PlayerState) []*discordgo.MessageEmbedField {
	return lo.Map(state.Playlist, func(item *music.Track, index int) *discordgo.MessageEmbedField {
		return &discordgo.MessageEmbedField{
			Name:  "> " + item.Title + " - " + formatTrackArtists(item) + " ``` " + secondsToMinutes(item.Duration) + " ```",
			Value: "",
		}
	})
}

func thumbnail(state music.PlayerState) string {
	if len(state.Playlist) <= 0 {
		return ""
	}
	return state.Playlist[0].Thumbnail
}

func formatTrackArtists(track *music.Track) string {
	return strings.Join(lo.Map(track.Artists, func(a music.Artist, i int) string { return a.Name }), ", ")
}

func secondsToMinutes(millis uint) string {
	inSeconds := millis / 1000
	minutes := inSeconds / 60
	seconds := inSeconds % 60
	return fmt.Sprintf("%d:%d", minutes, seconds)
}
