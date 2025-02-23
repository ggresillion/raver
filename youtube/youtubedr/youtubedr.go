package youtubedr

import (
	"encoding/json"
	"raver/youtube"

	ytdr "github.com/kkdai/youtube/v2"
)

type YoutubeDR struct{}

func NewYoutubeDRAdapter() youtube.YoutubeAdapter {
	return &YoutubeDR{}
}

func (y *YoutubeDR) GetVideoByID(videoID string) (*youtube.Video, error) {
	yt := ytdr.Client{}
	video, err := yt.GetVideo(videoID)
	if err != nil {
		return nil, err
	}

	format := video.Formats.WithAudioChannels()[0]

	v, _ := json.Marshal(format)
	println(string(v))

	url, err := yt.GetStreamURL(video, &format)
	if err != nil {
		return nil, err
	}
	println(url)

	return &youtube.Video{
		ID:       video.ID,
		URL:      url,
		Title:    video.Title,
		Channel:  video.Author,
		Duration: video.Duration,
		IsLive:   false,
	}, nil
}
