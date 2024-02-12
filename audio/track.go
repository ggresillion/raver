package audio

import (
	"io"
	"log"
	"time"
)

type TrackInfo struct {
	ID       string
	Title    string
	Artist   string
	Duration time.Duration
	Live     bool
}

type AudioStream struct {
	ProgressBytes int64
	TotalBytes    int64
	io.ReadCloser
}

func NewAudioStream(in io.ReadCloser, length int64) *AudioStream {
	s := &AudioStream{
		ReadCloser: in,
		TotalBytes: length,
	}
	log.Printf("stream[%p]: created a new stream", s)
	return s
}

func (s *AudioStream) Read(bytes []byte) (n int, err error) {
	n, err = s.ReadCloser.Read(bytes)
	s.ProgressBytes += int64(n)
	return
}

type Track struct {
	TrackInfo
	*AudioStream
}

func NewTrack(infos TrackInfo, stream *AudioStream) *Track {
	return &Track{infos, stream}
}
