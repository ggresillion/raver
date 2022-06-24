package bot

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/bwmarrin/discordgo"
)

var ErrNotNotInVoiceChannel = errors.New("bot is not in a voice channel")
var ErrNoVoiceConnection = errors.New("no voice connection")
var ErrCouldNotSetSpeaking = errors.New("couldn't set speaking")
var ErrCouldNotStartStream = fmt.Errorf("couldn't start stream")

type Stream struct {
	Progress chan time.Duration
	End      chan error
}

type Audio struct {
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

func NewAudio(guildID string, bot *Bot) *Audio {
	return &Audio{guildID: guildID, bot: bot, audioStatus: NotConnected, audioStatusChange: make(chan AudioStatus)}
}

func (b *Audio) LeaveChannel() {
	vc := b.voiceConnection
	if vc == nil {
		return
	}
	vc.Disconnect()
	vc.Close()
	b.voiceConnection = nil
}

func (b *Audio) JoinUserChannel(userID string) error {
	g, err := b.bot.session.State.Guild(b.guildID)
	if err != nil {
		return err
	}
	defer log.Printf("[%s] joinned voice channel", b.guildID)
	for _, v := range g.VoiceStates {
		if v.UserID == userID {
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

func (b *Audio) Play(url string, start time.Duration) (*Stream, error) {

	switch b.audioStatus {
	case NotConnected:
		return nil, ErrNotNotInVoiceChannel
	case IDLE:
		break
	}

	if b.voiceConnection == nil {
		return nil, ErrNoVoiceConnection
	}

	err := b.voiceConnection.Speaking(true)
	if err != nil {
		return nil, ErrCouldNotSetSpeaking
	}

	b.session, err = b.startStream(url, start)
	if err != nil {
		return nil, ErrCouldNotStartStream
	}

	return &Stream{
		Progress: b.session.progress,
		End:      b.session.end,
	}, nil
}

func (b *Audio) Status() AudioStatus {
	return b.audioStatus
}

func (b *Audio) StatusChange() chan AudioStatus {
	return b.audioStatusChange
}

func (b *Audio) setStatus(s AudioStatus) {
	b.audioStatus = s
	b.audioStatusChange <- s
}
