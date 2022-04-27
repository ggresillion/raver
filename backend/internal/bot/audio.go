package bot

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fiwippi/dca"
)

type BotAudio struct {
	guildID           string
	bot               *Bot
	voiceStates       []*discordgo.VoiceState
	voiceConnection   *discordgo.VoiceConnection
	audioStatus       AudioStatus
	audioStatusChange chan AudioStatus
	session           *StreamingSession
}

type StreamingSession struct {
	url              string
	end              chan error
	encodingSession  *dca.EncodeSession
	streamingSession *dca.StreamingSession
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

func (b *BotAudio) Play(url string) (chan error, error) {
	switch b.audioStatus {
	case NotConnected:
		return nil, errors.New("bot is not in a voice channel")
	case Playing:
		b.Stop()
	}

	if b.voiceConnection == nil {
		return nil, errors.New("no voice connection")
	}

	err := b.voiceConnection.Speaking(true)
	if err != nil {
		return nil, errors.New("couldn't set speaking")
	}

	b.session, err = b.getStream(url, time.Duration(0))
	if err != nil {
		return nil, fmt.Errorf("couldn't start stream: %w", err)
	}

	return b.session.end, nil
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
	b.Stop()

	var err error
	b.session, err = b.getStream(url, t)
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

func (b *BotAudio) getStream(url string, start time.Duration) (*StreamingSession, error) {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"
	options.StartTime = int(start.Seconds())

	session := &StreamingSession{
		url:              url,
		end:              make(chan error),
		encodingSession:  nil,
		streamingSession: nil,
	}

	var err error
	session.encodingSession, err = dca.EncodeFile(url, options)
	if err != nil {
		return nil, err
	}

	streamEnd := make(chan error)

	session.streamingSession = dca.NewStream(session.encodingSession, b.voiceConnection, streamEnd)

	// Listen for the stream ending
	go func(end chan error) {
		err := <-streamEnd
		b.setStatus(IDLE)
		session.encodingSession.Cleanup()
		b.voiceConnection.Speaking(false)
		if err != nil && err != io.EOF {
			end <- err
		} else {
			end <- nil
		}
		b.session = nil
	}(session.end)

	b.setStatus(Playing)

	return session, nil
}
