package bot

import (
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"

	"github.com/ggresillion/dca"
)

type StreamingSession struct {
	sync.Mutex
	End              chan error
	Progress         chan time.Duration
	url              string
	options          *dca.EncodeOptions
	encodingSession  *dca.EncodeSession
	streamingSession *dca.StreamingSession
	streamEnd        chan error
}

func newStreamingSession(url string, options *dca.EncodeOptions) *StreamingSession {
	return &StreamingSession{
		End:              make(chan error),
		url:              url,
		Progress:         make(chan time.Duration),
		options:          options,
		encodingSession:  nil,
		streamingSession: nil,
		streamEnd:        make(chan error),
	}
}

// startStream starts a new stream given an url, it must be an audio or video stream that can be decoded by ffmpeg.
// Optional parameter start allow the start the stream at a specific time.
func (a *Audio) startStream(url string, start time.Duration) (*StreamingSession, error) {
	// Encoding options
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"
	options.StartTime = int(start.Seconds())

	// Ensure to stop the current stream if any
	a.Stop()

	a.session = newStreamingSession(url, options)

	// Configure the encoding session
	var err error
	a.session.encodingSession, err = dca.EncodeFile(url, a.session.options)
	if err != nil {
		return nil, err
	}

	// Start the stream
	a.session.streamingSession = dca.NewStream(a.session.encodingSession, a.voiceConnection, a.session.streamEnd)

	// Listen for the stream ending
	go a.listenToStreamEnding()

	// Updates the stream Progress
	go a.updateStreamProgress()

	a.setStatus(Playing)

	return a.session, nil
}

// updateStreamProgress update the stream current progress by fetching periodically from the streaming session.
// Stops when the stream is ended.
func (a *Audio) updateStreamProgress() {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		<-ticker.C
		if a.session == nil {
			return
		}
		a.session.Lock()
		startTime := time.Duration(a.session.options.StartTime) * time.Second
		a.session.Progress <- startTime + a.session.streamingSession.PlaybackPosition()
		a.session.Unlock()
	}
}

// listenToStreamEnding trigger cleanup when the stream has ended.
func (a *Audio) listenToStreamEnding() {
	err := <-a.session.streamEnd
	a.session.Lock()
	a.setStatus(IDLE)
	a.session.streamingSession.SetPaused(true)
	a.session.encodingSession.Cleanup()
	err = a.voiceConnection.Speaking(false)
	if err != nil {
		log.Errorf("[%s] (listenToStreamEnding) could not set speaking status: %e", a.guildID, err)
	}

	if err != nil && err != io.EOF {
		a.session.End <- err
	} else {
		a.session.End <- nil
	}
	a.session = nil
}

// Pause the stream.
func (a *Audio) Pause() {
	if a.session == nil || a.audioStatus != Playing {
		return
	}
	a.session.streamingSession.SetPaused(true)
	a.setStatus(Paused)
}

// UnPause the stream.
func (a *Audio) UnPause() {
	if a.session == nil || a.audioStatus != Paused {
		return
	}
	a.session.streamingSession.SetPaused(false)
	a.setStatus(Playing)
}

// Stop the stream.
func (a *Audio) Stop() {
	if a.session == nil {
		return
	}
	a.session.streamEnd <- nil
}
