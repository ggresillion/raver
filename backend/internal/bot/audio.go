package bot

import (
	"errors"
	"io"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
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
	end              chan error
	stop             chan error
	pause            chan bool
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

	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	b.session = &StreamingSession{}

	b.session.encodingSession, err = dca.EncodeFile(url, options)
	if err != nil {
		return nil, err
	}

	b.session.end = make(chan error)
	b.session.stop = make(chan error)
	b.session.pause = make(chan bool)

	b.setStatus(Playing)

	streamEnd := make(chan error)

	b.session.streamingSession = dca.NewStream(b.session.encodingSession, b.voiceConnection, streamEnd)

	// Listen for pause or stop messages
	go func(pause chan bool, stop chan error) {
		for {
			select {
			case p := <-pause:
				b.session.streamingSession.SetPaused(p)
				if p {
					b.setStatus(Paused)
				} else {
					b.setStatus(Playing)
				}
			case <-stop:
				b.session.encodingSession.Cleanup()
			}
		}
	}(b.session.pause, b.session.stop)

	// Listen for the stream ending
	go func(end chan error) {
		err := <-streamEnd
		b.session.encodingSession.Cleanup()
		b.voiceConnection.Speaking(false)
		b.setStatus(IDLE)
		b.session = nil
		if err != nil && err != io.EOF {
			end <- err
		} else {
			end <- nil
		}
	}(b.session.end)

	return b.session.end, nil
}

func (b *BotAudio) Pause() {
	if b.session.pause == nil || b.audioStatus != Playing {
		return
	}
	b.session.pause <- true
}

func (b *BotAudio) UnPause() {
	if b.session.pause == nil || b.audioStatus != Paused {
		return
	}
	b.session.pause <- false
}

func (b *BotAudio) Stop() {
	b.session.stop <- nil
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
