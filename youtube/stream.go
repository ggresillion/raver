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

func GetStreamFromYoutube(videoID string) (*audio.AudioStream, error) {
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return nil, err
	}

	log.Printf("Got video info %q", video.ID)

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
	log.Println("Got video stream")

	audioStream, err := extractOpus(rs)
	if err != nil {
		return nil, err
	}

	return audioStream, nil
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
