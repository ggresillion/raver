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

func GetPlayableTrackFromYoutube(videoID string) (*audio.Track, error) {
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return nil, err
	}

	log.Printf("Youtube: got video info %q", video.ID)

	formats := video.Formats.FindByItag(251) // only get videos with audio

	url, err := client.GetStreamURL(video, formats)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	rs := httprs.NewHttpReadSeeker(resp)
	log.Println("Youtube: got webm video stream")

	audioStream, err := extractOpus(rs)
	if err != nil {
		return nil, err
	}
	log.Printf("Youtube: converting to opus stream %p", audioStream)

	return audio.NewTrack(audio.TrackInfo{video.ID, video.Title, video.Author, video.Duration, false}, audioStream), nil
}

func extractOpus(stream io.ReadSeeker) (*audio.AudioStream, error) {
	audioStream := audio.NewAudioStream()
	var w webm.WebM
	wr, err := webm.Parse(stream, &w)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			el := <-wr.Chan
			audioStream.In <- el.Data
		}
	}()

	return audioStream, nil
}
