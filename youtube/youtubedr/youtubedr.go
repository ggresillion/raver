package youtubedr

import (
	"errors"
	"log/slog"
	"raver/youtube"
	"strings"

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

	var format *ytdr.Format
	for _, f := range video.Formats.WithAudioChannels() {
		if strings.Contains(f.MimeType, "opus") {
			format = &f
		}
	}

	if format == nil {
		return nil, errors.New("no opus format found")
	}

	url, err := yt.GetStreamURL(video, format)
	if err != nil {
		return nil, err
	}

	slog.Info("[youtubedr] got stream url", "format", format.MimeType)

	return &youtube.Video{
		ID:       video.ID,
		URL:      url,
		Title:    video.Title,
		Channel:  video.Author,
		Duration: video.Duration,
		IsLive:   false,
	}, nil
}
