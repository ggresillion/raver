package display

import (
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

func ClearMessages(s *discordgo.Session, channelID string) error {
	messages, err := s.ChannelMessages(channelID, 100, "", "", "")
	if err != nil {
		return err
	}

	return s.ChannelMessagesBulkDelete(
		channelID,
		lo.Map(messages, func(m *discordgo.Message, i int) string { return m.ID }))
}
