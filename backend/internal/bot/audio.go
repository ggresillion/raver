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

func (b *BotAudio) Play(url string) (chan error, chan time.Duration, error) {

	if b.session != nil {
		return nil, nil, errors.New("bot is already playing something")
	}

	switch b.audioStatus {
	case NotConnected:
		return nil, nil, errors.New("bot is not in a voice channel")
	case IDLE:
		break
	case Playing:
		return nil, nil, errors.New("bot is already playing something")
	}

	if b.voiceConnection == nil {
		return nil, nil, errors.New("no voice connection")
	}

	err := b.voiceConnection.Speaking(true)
	if err != nil {
		return nil, nil, errors.New("couldn't set speaking")
	}

	b.session, err = b.startStream(url, time.Duration(0))
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't start stream: %w", err)
	}

	return b.session.end, b.session.progress, nil
}

func (b *BotAudio) Pause() {
	if b.session == nil || b.audioStatus != Playing {
		return
	}
	b.session.streamingSession.SetPaused(true)
	b.setStatus(Paused)
}

func (b *BotAudio) UnPause() {
	if b.session == nil || b.audioStatus != Paused {
		return
	}
	b.session.streamingSession.SetPaused(false)
	b.setStatus(Playing)
}

func (b *BotAudio) SetTime(t time.Duration) error {
	url := b.session.url
	end := b.session.end
	b.session.swappingStream = true
	b.Stop()

	var err error
	b.session, err = b.startStream(url, t)
	b.session.end = end
	if err != nil {
		return fmt.Errorf("couldn't start stream: %w", err)
	}
	return nil
}

func (b *BotAudio) Stop() {
	b.session.encodingSession.Cleanup()
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
