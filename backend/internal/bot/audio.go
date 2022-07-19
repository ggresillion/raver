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

func (a *Audio) LeaveChannel() {
	vc := a.voiceConnection
	if vc == nil {
		return
	}
	vc.Disconnect()
	vc.Close()
	a.voiceConnection = nil
}

func (a *Audio) JoinUserChannel(userID string) error {
	g, err := a.bot.session.State.Guild(a.guildID)
	if err != nil {
		return err
	}
	defer log.Printf("[%s] joinned voice channel", a.guildID)
	for _, v := range g.VoiceStates {
		if v.UserID == userID {
			vc, err := a.bot.session.ChannelVoiceJoin(a.guildID, v.ChannelID, false, true)
			if err != nil {
				return err
			}
			a.voiceConnection = vc

			if a.audioStatus == NotConnected {
				a.setStatus(IDLE)
			}

			return nil
		}
	}
	return errors.New("user not in a voice channel")
}

func (a *Audio) Play(url string, start time.Duration) (*StreamingSession, error) {

	switch a.audioStatus {
	case NotConnected:
		return nil, ErrNotNotInVoiceChannel
	case IDLE:
		break
	}

	if a.voiceConnection == nil {
		return nil, ErrNoVoiceConnection
	}

	err := a.voiceConnection.Speaking(true)
	if err != nil {
		return nil, ErrCouldNotSetSpeaking
	}

	a.session, err = a.startStream(url, start)
	if err != nil {
		return nil, ErrCouldNotStartStream
	}

	return a.session, nil
}

func (a *Audio) Status() AudioStatus {
	return a.audioStatus
}

func (a *Audio) StatusChange() chan AudioStatus {
	return a.audioStatusChange
}

func (a *Audio) setStatus(s AudioStatus) {
	a.audioStatus = s
	a.audioStatusChange <- s
}
