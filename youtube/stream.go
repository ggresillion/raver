package youtube

import (
	"errors"
	"io"
	"log"
	"net/http"

	"raver/audio"

	"github.com/ebml-go/webm"
	"github.com/jfbus/httprs"
	"github.com/kkdai/youtube/v2"
)

// GetPlayableTrackFromYoutube returns a audio.Track from a given videoID
func GetPlayableTrackFromYoutube(videoID string) (*audio.Track, error) {
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return nil, err
	}

	log.Printf("youtube: got video info %q", video.ID)

	formats := video.Formats.Quality("251") // only get videos with audio

	if len(formats) < 1 {
		return nil, errors.New("youtube: could not find suitable audio format")
	}

	url, err := client.GetStreamURL(video, &formats[0])
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	rs := httprs.NewHttpReadSeeker(resp)
	log.Println("youtube: got webm video stream")

	audioStream, err := extractOpus(rs, resp.ContentLength)
	if err != nil {
		return nil, err
	}
	log.Println("youtube: converting to opus stream")

	return audio.NewTrack(audio.TrackInfo{ID: video.ID, Title: video.Title, Artist: video.Author, Duration: video.Duration, Live: false}, audioStream), nil
}

// extractOpus reads the incoming stream, parses it as a webm container and extract opus stream.
func extractOpus(stream io.ReadSeeker, length int64) (*audio.AudioStream, error) {
	var w webm.WebM
	wr, err := webm.Parse(stream, &w)
	if err != nil {
		return nil, err
	}

	in := NewYTReadCloser(wr)
	log.Printf("youtube[%p]: created new input stream", &in)
	audioStream := audio.NewAudioStream(in, length)

	return audioStream, nil
}

type YTReadCloser struct {
	wr *webm.Reader
}

func NewYTReadCloser(wr *webm.Reader) *YTReadCloser {
	return &YTReadCloser{wr: wr}
}

func (r *YTReadCloser) Read(bytes []byte) (n int, err error) {
	packet, ok := <-r.wr.Chan
	if !ok {
		log.Printf("youtube[%p]: closed stream", &r)
		return 0, io.EOF
	}
	if len(packet.Data) == 0 {
		log.Printf("youtube[%p]: end of input stream", &r)
		r.wr.Shutdown()
		return 0, io.EOF
	}
	copy(bytes, packet.Data)
	return len(packet.Data), nil
}

func (r *YTReadCloser) Close() (err error) {
	log.Printf("youtube[%p]: closing stream", &r)
	r.wr.Shutdown()
	return
}
