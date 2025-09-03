package goytdlp

import (
	"context"
	"fmt"
	"raver/youtube"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

type YoutubeAdapter struct{}

func NewYoutubeAdapter() YoutubeAdapter {
	ytdlp.MustInstall(context.TODO(), nil)
	return YoutubeAdapter{}
}

func (y YoutubeAdapter) GetVideoByID(videoID string) (*youtube.Video, error) {
	ctx := context.TODO()

	cmd := ytdlp.New().
		Format("bestaudio[ext=webm][acodec=opus]").
		NoProgress().
		PrintJSON()

	res, err := cmd.Run(ctx, fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID))
	if err != nil {
		return nil, err
	}

	info, err := res.GetExtractedInfo()
	if err != nil {
		return nil, err
	}

	video := info[0]

	return &youtube.Video{
		ID:       video.ID,
		URL:      *video.URL,
		Title:    *video.Title,
		Channel:  *video.Channel,
		Duration: time.Duration(*video.Duration * float64(time.Second)),
	}, nil
}
