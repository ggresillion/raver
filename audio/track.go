package audio

import (
	"io"
	"log"
	"sync"
	"time"
)

const audioBufferSize = 10

type AudioBuffer struct {
	c        chan []byte
	closed   bool
	closeMux sync.Mutex
}

type TrackInfo struct {
	ID       string
	Title    string
	Artist   string
	Duration time.Duration
	Live     bool
}

type AudioStream struct {
	Progress time.Duration
	buffer   *AudioBuffer
	playing  bool
}

func NewAudioStream() *AudioStream {
	s := &AudioStream{
		buffer: &AudioBuffer{
			c: make(chan []byte, audioBufferSize),
		},
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
	s.buffer.closeMux.Lock()
	defer s.buffer.closeMux.Unlock()

	if s.buffer.closed {
		return io.EOF
	}

	s.buffer.c <- bytes
	return nil
}

func (s *AudioStream) Read() ([]byte, error) {
	data, ok := <-s.buffer.c
	if !ok || len(data) == 0 {
		return nil, io.EOF
	}
	return data, nil
}

func (s *AudioStream) Stop() {
	s.buffer.closeMux.Lock()
	defer s.buffer.closeMux.Unlock()

	if !s.buffer.closed {
		close(s.buffer.c)
		s.buffer.closed = true
		log.Printf("stream[%p]: stopped stream", s)
	}
}
