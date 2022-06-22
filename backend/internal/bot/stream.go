package bot

import (
	"io"
	"time"

	"github.com/ggresillion/dca"
)

type StreamingSession struct {
	url              string
	end              chan error
	naturalEnd       chan error
	options          *dca.EncodeOptions
	progress         chan time.Duration
	encodingSession  *dca.EncodeSession
	streamingSession *dca.StreamingSession
	swappingStream   bool
}

func (b *BotAudio) startStream(url string, start time.Duration) (*StreamingSession, error) {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"
	options.StartTime = int(start.Seconds())

	b.session = &StreamingSession{
		url:              url,
		end:              make(chan error),
		progress:         make(chan time.Duration),
		options:          options,
		encodingSession:  nil,
		streamingSession: nil,
	}

	var err error
	b.session.encodingSession, err = dca.EncodeFile(url, b.session.options)
	if err != nil {
		return nil, err
	}

	naturalEnd := make(chan error)
	stopProgress := make(chan bool)

	b.session.streamingSession = dca.NewStream(b.session.encodingSession, b.voiceConnection, naturalEnd)

	// Listen for the stream ending
	go func(end chan error) {
		for {
			err := <-naturalEnd
			stopProgress <- true
			b.setStatus(IDLE)
			if b.session != nil {
				b.session.encodingSession.Cleanup()
			}
			b.voiceConnection.Speaking(false)
			if err != nil && err != io.EOF {
				end <- err
			} else {
				end <- nil
			}
			b.session = nil
			return
		}
	}(b.session.end)

	// Updates the stream progress
	go b.session.updateStreamProgress(stopProgress)

	b.setStatus(Playing)

	return b.session, nil
}

func (s *StreamingSession) updateStreamProgress(stop chan bool) {
	for {
		select {
		case <-stop:
			return
		default:
			startTime := time.Duration(s.options.StartTime) * time.Second
			s.progress <- startTime + s.streamingSession.PlaybackPosition()
			time.Sleep(time.Second)
		}
	}
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

func (b *BotAudio) Stop() {
	if b.session == nil {
		return
	}
	b.audioStatus = IDLE
	b.session.encodingSession.Cleanup()
}
