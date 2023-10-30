package audio

type AudioStream struct {
	// Input audio stream
	In chan []byte
	// Output audio stream
	Out     chan []byte
	End     chan error
	queue   chan []byte
	pause   chan bool
	playing bool
}

func NewAudioStream() *AudioStream {
	return &AudioStream{
		In:      make(chan []byte),
		queue:   make(chan []byte),
		End:     make(chan error),
		pause:   make(chan bool),
		playing: false,
	}
}

func (s *AudioStream) Bind(c chan []byte) {
	s.Out = c
}

func (s *AudioStream) Play() {
	if s.playing {
		s.pause <- true
		return
	}
	s.playing = true
	go func() {
	play:
		for {
			select {
			case <-s.pause:
				goto pause
			case <-s.End:
				s.playing = false
				return
			case data, more := <-s.In:
				if !more {
					s.playing = false
					s.End <- nil
					return
				}
				s.Out <- data
			}
		}
	pause:
		<-s.pause
		goto play
	}()
}

func (s *AudioStream) Stop() {
	s.End <- nil
}

func (s *AudioStream) Pause() {
	s.pause <- true
}

func (s *AudioStream) BlockUntilStop() {
	<-s.End
}
