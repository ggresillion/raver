package audio

import (
	"io"
	"log"
	"math"
	"time"
)

const (
	maxBufferedBytes = 1000
	maxBytesPerFrame = 960
)

type TrackInfo struct {
	ID       string
	Title    string
	Artist   string
	Duration time.Duration
	Live     bool
}

type AudioStream struct {
	io.ReadCloser
	ProgressBytes int64
	TotalBytes    int64
	buffer        chan []byte
}

func NewAudioStream(in io.ReadCloser, length int64) *AudioStream {
	maxBufferedFrames := int(math.Ceil(float64(maxBufferedBytes) / float64(maxBytesPerFrame)))
	s := &AudioStream{
		ReadCloser: in,
		TotalBytes: length,
		buffer:     make(chan []byte, maxBufferedFrames),
	}
	log.Printf("stream[%p]: created a new stream (max_buffered_bytes: %d, max_bytes_per_frame: %d, max_buffered_frames: %d", s, maxBufferedBytes, maxBytesPerFrame, maxBufferedFrames)
	go func() {
		for {
			buf := make([]byte, maxBufferedBytes)
			n, err := in.Read(buf)
			if err != nil {
				close(s.buffer)
				break
			}
			s.buffer <- buf[:n]
		}
	}()
	return s
}

func (s *AudioStream) Read(bytes []byte) (n int, err error) {
	buf, ok := <-s.buffer
	if !ok {
		err = io.EOF
		return
	}
	n = copy(bytes, buf)
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
