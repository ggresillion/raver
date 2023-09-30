package audio

type AudioStream struct {
	// Input audio stream
	In chan []byte
	// Output audio stream
	Out     chan []byte
	stop    chan error
	queue   chan []byte
	pause   chan bool
	playing bool
}

func NewAudioStream() *AudioStream {
	return &AudioStream{
		In:      make(chan []byte),
		Out:     make(chan []byte),
		queue:   make(chan []byte),
		stop:    make(chan error),
		pause:   make(chan bool),
		playing: false,
	}
}

func (s *AudioStream) Play() {
	s.playing = true
	go func() {
	play:
		for {
			select {
			case <-s.pause:
				goto pause
			case <-s.stop:
				s.playing = false
				return
			case s.Out <- <-s.In:
			}
		}
	pause:
		<-s.pause
		goto play
	}()
}

func (s *AudioStream) Stop() {
	s.stop <- nil
}

func (s *AudioStream) Pause() {
	s.pause <- true
}

func (s *AudioStream) Resume() {
	s.pause <- false
}

func (s *AudioStream) BlockUntilStop() {
	<-s.stop
}
