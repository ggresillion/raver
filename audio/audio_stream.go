package audio

import (
	"log"
)

const (
	bufferSize = 255
)

type AudioStream struct {
	In          chan []byte
	Out         chan []byte
	buffer      chan []byte
	end         chan error
	endChannels []chan struct{}
	queue       chan []byte
	pause       chan bool
	playing     bool
}

func NewAudioStream(in chan []byte) *AudioStream {
	s := &AudioStream{
		In:          in,
		Out:         make(chan []byte),
		buffer:      make(chan []byte, bufferSize),
		end:         make(chan error),
		endChannels: make([]chan struct{}, 0),
		queue:       make(chan []byte),
		pause:       make(chan bool),
		playing:     false,
	}
	// buffering all input frames
	go func() {
		for {
			select {
			case data, more := <-s.In:
				if !more {
					log.Printf("stream[%p]: end of stream", s)
					close(s.buffer)
					return
				}
				s.buffer <- data
			case <-s.OnStop():
				log.Printf("stream[%p]: closing buffer", s)
				close(s.buffer)
				return
			}
		}
	}()
	log.Printf("stream[%p]: created a new stream", s)
	return s
}

func (s *AudioStream) Play() {
	if s.playing {
		s.pause <- true
		log.Printf("stream[%p]: pausing stream", s)
		return
	}
	s.playing = true
	log.Printf("stream[%p]: playing stream", s)
	// main loop
	go func() {
		defer func() {
			s.playing = false
			for _, ch := range s.endChannels {
				close(ch)
			}
			close(s.Out)
			close(s.pause)
			close(s.queue)
			close(s.end)
			log.Printf("stream[%p]: closed stream", s)
		}()
	play:
		for {
			select {
			case <-s.pause:
				log.Printf("stream[%p]: pausing stream", s)
				goto pause
			case <-s.end:
				log.Printf("stream[%p]: stopping stream", s)
				return
			case data, ok := <-s.buffer:
				if !ok {
					log.Printf("stream[%p]: no more data in buffer, closing", s)
					return
				}
				s.Out <- data
			}
		}
	pause:
		<-s.pause
		log.Printf("stream[%p]: resuming stream", s)
		goto play
	}()
}

func (s *AudioStream) Stop() {
	if !s.playing {
		log.Printf("stream[%p]: not playing", s)
		return
	}
	go func() { s.end <- nil }()
}

func (s *AudioStream) OnStop() <-chan struct{} {
	ch := make(chan struct{})
	s.endChannels = append(s.endChannels, ch)
	return ch
}

func (s *AudioStream) Pause() {
	s.pause <- true
}
