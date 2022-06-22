package bot

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

type BotAudio struct {
	// Guild ID
	guildID string
	// Bot
	bot *Bot
	// Voice states of all the users in the guild (used in voice channel join)
	voiceStates []*discordgo.VoiceState
	// Audio status of the bot
	audioStatus       AudioStatus
	audioStatusChange chan AudioStatus
	// Voice connection for the bot
	voiceConnection *discordgo.VoiceConnection
	// Infos about the current stream playing
	session *StreamingSession
}

func NewBotAudio(guildID string, bot *Bot) *BotAudio {
	return &BotAudio{guildID: guildID, bot: bot, audioStatus: NotConnected, audioStatusChange: make(chan AudioStatus)}
}

func (b *BotAudio) LeaveChannel() {
	vc := b.voiceConnection
	if vc == nil {
		return
	}
	vc.Disconnect()
	vc.Close()
	b.voiceConnection = nil
}

func (b *BotAudio) JoinUserChannel(userID string) error {
	g, err := b.bot.session.State.Guild(b.guildID)
	if err != nil {
		return err
	}
	for _, v := range g.VoiceStates {
		if v.UserID == userID {
			defer log.Printf("joinned voice channel for guildID %s and userID %s", b.guildID, userID)
			vc, err := b.bot.session.ChannelVoiceJoin(b.guildID, v.ChannelID, false, true)
			if err != nil {
				return err
			}
			b.voiceConnection = vc

			if b.audioStatus == NotConnected {
				b.setStatus(IDLE)
			}

			return nil
		}
	}
	return errors.New("user not in a voice channel")
}

func (b *BotAudio) Play(url string, start time.Duration) (chan error, chan time.Duration, error) {

	switch b.audioStatus {
	case NotConnected:
		return nil, nil, errors.New("bot is not in a voice channel")
	case IDLE:
		break
	}

	if b.voiceConnection == nil {
		return nil, nil, errors.New("no voice connection")
	}

	err := b.voiceConnection.Speaking(true)
	if err != nil {
		return nil, nil, errors.New("couldn't set speaking")
	}

	b.session, err = b.startStream(url, start)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't start stream: %w", err)
	}

	return b.session.end, b.session.progress, nil
}

func (b *BotAudio) Status() AudioStatus {
	return b.audioStatus
}

func (b *BotAudio) StatusChange() chan AudioStatus {
	return b.audioStatusChange
}

func (b *BotAudio) setStatus(s AudioStatus) {
	b.audioStatus = s
	b.audioStatusChange <- s
}
