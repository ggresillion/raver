package bot

import (
	"errors"
	"io"
	"log"

	"github.com/jonas747/dca"
)

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
	g, err := b.bot.session.State.Guild(b.guildId)
	if err != nil {
		return err
	}
	for _, v := range g.VoiceStates {
		if v.UserID == userID {
			defer log.Printf("joinned voice channel for guildID %s and userID %s", b.guildId, userID)
			vc, err := b.bot.session.ChannelVoiceJoin(b.guildId, v.ChannelID, false, true)
			if err != nil {
				return err
			}
			b.voiceConnection = vc
			b.audioStatus = IDLE
			return nil
		}
	}
	return errors.New("user not in a voice channel")
}

func (b *BotAudio) Play(url string) error {

	v := b.voiceConnection
	if v == nil {
		return errors.New("bot is not in voice channel")
	}

	err := v.Speaking(true)
	if err != nil {
		return errors.New("couldn't set speaking")
	}

	defer func() {
		v.Speaking(false)
	}()

	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	encodingSession, err := dca.EncodeFile(url, options)
	if err != nil {
		return err
	}

	done := make(chan error)
	pause := make(chan bool)

	go func() {
		defer encodingSession.Cleanup()
		defer close(pause)
		defer close(done)

		ss := dca.NewStream(encodingSession, v, done)

		go func() {
			for {
				select {
				case p := <-pause:
					ss.SetPaused(p)
					if p {
						b.audioStatus = Paused
					} else {
						b.audioStatus = Playing
					}
				case <-done:
					encodingSession.Stop()
					b.audioStatus = IDLE
					b.pause = nil
					b.stop = nil
					return
				}
			}
		}()

		err = <-done
		if err != nil && err != io.EOF {
			log.Printf("error playing audio stream: %s", err)
		}
	}()

	b.audioStatus = Playing
	b.pause = pause
	b.stop = done

	return nil
}

func (b *BotAudio) Pause() {
	if b.pause == nil {
		return
	}

	if b.audioStatus == Playing {
		b.pause <- true
	} else if b.audioStatus == Paused {
		b.pause <- false
	}
}

func (b *BotAudio) Stop() {
	b.stop <- nil
}

func (b *BotAudio) Status() AudioStatus {
	return b.audioStatus
}
