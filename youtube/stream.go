package youtube

import (
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

	audioStream, err := extractOpus(rs)
	if err != nil {
		return nil, err
	}
	log.Println("youtube: converting to opus stream")

	return audio.NewTrack(audio.TrackInfo{ID: video.ID, Title: video.Title, Artist: video.Author, Duration: video.Duration, Live: false}, audioStream), nil
}

// extractOpus reads the incoming stream, parses it as a webm container and extract opus stream.
func extractOpus(stream io.ReadSeeker) (*audio.AudioStream, error) {
	audioStream := audio.NewAudioStream()
	var w webm.WebM
	wr, err := webm.Parse(stream, &w)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			packet, ok := <-wr.Chan
			if len(packet.Data) == 0 || !ok {
				log.Printf("youtube[%p]: end of input stream", audioStream)
				audioStream.End()
			}
			err := audioStream.Write(packet.Data)
			if err == io.EOF {
				wr.Shutdown()
				return
			}
		}
	}()

	return audioStream, nil
}
