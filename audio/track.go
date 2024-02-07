package audio

import (
	"io"
	"log"
	"time"
)

const audioBufferSize = 200
const endOfFile = "EOF"

type TrackInfo struct {
	ID       string
	Title    string
	Artist   string
	Duration time.Duration
	Live     bool
}

type AudioStream struct {
	buffer  chan []byte
	playing bool
	stopped bool
	close   chan struct{}
}

func NewAudioStream() *AudioStream {
	s := &AudioStream{
		buffer:  make(chan []byte, audioBufferSize),
		playing: false,
		stopped: false,
	}
	log.Printf("stream[%p]: created a new stream", s)
	return s
}

type Track struct {
	TrackInfo
	*AudioStream
}

func NewTrack(infos TrackInfo, stream *AudioStream) *Track {
	return &Track{infos, stream}
}

func (s *AudioStream) Write(bytes []byte) error {
	if s.stopped {
		return io.EOF
	}
	s.buffer <- bytes
	return nil
}

func (s *AudioStream) Read() ([]byte, error) {
	if s.stopped {
		return nil, io.EOF
	}
	data := <-s.buffer
	if len(data) == 0 {
		return nil, io.EOF
	}
	return data, nil
}

func (s *AudioStream) Play() {
	s.playing = true
	log.Printf("stream[%p]: playing stream", s)
}

func (s *AudioStream) Pause() {
	s.playing = false
	log.Printf("stream[%p]: paused stream", s)
}

func (s *AudioStream) Resume() {
	s.playing = true
	log.Printf("stream[%p]: paused stream", s)
}

func (s *AudioStream) Stop() {
	s.stopped = true
	close(s.buffer)
	log.Printf("stream[%p]: stopped stream", s)
}

func (s *AudioStream) Close() {
	s.buffer <- []byte{}
}
