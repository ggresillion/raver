package bot

import (
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) GetVoiceConnection(guildID string) (*discordgo.VoiceConnection, error) {
	vc := b.voiceConnections[guildID]
	if vc != nil {
		return vc, nil
	}
	return nil, fmt.Errorf("bot is not connected in a voice channel for guildID %s", guildID)
}

func (b *Bot) LeaveChannel(guildID string) {
	c := b.session.VoiceConnections[guildID]
	if c == nil {
		return
	}
	c.Disconnect()
	c.Close()
}

func (b *Bot) JoinUserChannel(guildID, userID string) (*discordgo.VoiceConnection, error) {
	g, err := b.session.State.Guild(guildID)
	if err != nil {
		return nil, err
	}
	for _, v := range g.VoiceStates {
		if v.UserID == userID {
			defer log.Printf("joinned voice channel for guildID %s and userID %s", guildID, userID)
			vc, err := b.session.ChannelVoiceJoin(guildID, v.ChannelID, false, true)
			if err != nil {
				return nil, err
			}
			b.voiceConnections[guildID] = vc
			return vc, nil
		}
	}
	return nil, errors.New("user not in a voice channel")
}
