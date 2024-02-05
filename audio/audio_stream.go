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
					log.Println("stream: end of stream")
					s.end <- nil
					return
				}
				s.buffer <- data
			case <-s.OnStop():
				log.Println("stream: closing buffer")
				close(s.buffer)
				return
			}
		}
	}()
	return s
}

func (s *AudioStream) Play() {
	if s.playing {
		s.pause <- true
		log.Println("stream: pausing stream")
		return
	}
	s.playing = true
	log.Println("stream: playing stream")
	// main loop
	go func() {
	play:
		for {
			select {
			case <-s.pause:
				log.Println("stream: pausing stream")
				goto pause
			case <-s.end:
				log.Println("stream: starting closing sequence")
				s.playing = false
				for _, ch := range s.endChannels {
					close(ch)
				}
				close(s.Out)
				close(s.pause)
				close(s.queue)
				close(s.end)
				log.Println("stream: closed stream")
				return
			case data := <-s.buffer:
				s.Out <- data
			}
		}
	pause:
		<-s.pause
		log.Println("stream: resuming stream")
		goto play
	}()
}

func (s *AudioStream) Stop() {
	if !s.playing {
		return
	}
	s.end <- nil
	log.Println("stream: sent end signal")
}

func (s *AudioStream) OnStop() <-chan struct{} {
	ch := make(chan struct{})
	s.endChannels = append(s.endChannels, ch)
	return ch
}

func (s *AudioStream) Pause() {
	s.pause <- true
}
