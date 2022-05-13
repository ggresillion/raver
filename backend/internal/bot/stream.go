package bot

import (
	"io"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ggresillion/dca"
)

type StreamingSession struct {
	url              string
	end              chan error
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
		encodingSession:  nil,
		streamingSession: nil,
	}

	var err error
	b.session.encodingSession, err = dca.EncodeFile(url, options)
	if err != nil {
		return nil, err
	}

	streamEnd := make(chan error)
	stopProgress := make(chan bool)

	b.session.streamingSession = dca.NewStream(b.session.encodingSession, b.voiceConnection, streamEnd)

	// Listen for the stream ending
	go func(end chan error) {

		err := <-streamEnd
		stopProgress <- true
		if b.session.swappingStream {
			return
		}
		b.setStatus(IDLE)
		b.session.encodingSession.Cleanup()
		b.voiceConnection.Speaking(false)
		if err != nil && err != io.EOF {
			end <- err
		} else {
			end <- nil
		}
		log.Debug(b.session.encodingSession.FFMPEGMessages())
		b.session = nil
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
			close(s.progress)
			return
		default:
			s.progress <- s.streamingSession.PlaybackPosition()
			time.Sleep(time.Second)
		}
	}
}
