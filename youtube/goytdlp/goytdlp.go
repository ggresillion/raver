package goytdlp

import (
	"context"
	"fmt"
	"log/slog"
	"raver/youtube"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

func init() {
	ytdlp.MustInstall(context.Background(), nil)
}

type YoutubeAdapter struct{}

func NewYoutubeAdapter() YoutubeAdapter {
	return YoutubeAdapter{}
}

func (y YoutubeAdapter) GetVideoByID(videoID string) (*youtube.Video, error) {
	ctx := context.TODO()

	ytdlp.Install(ctx, nil)

	cmd := ytdlp.New().
		Format("bestaudio[ext=webm][acodec=opus]").
		NoProgress().
		SkipDownload().
		PrintJSON()

	slog.Info("[ytdlp] getting info for video " + videoID + "...")

	res, err := cmd.Run(ctx, fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID))
	if err != nil {
		return nil, err
	}

	info, err := res.GetExtractedInfo()
	if err != nil {
		return nil, err
	}

	slog.Info("[ytdlp] got info for video " + videoID)

	video := info[0]

	return &youtube.Video{
		ID:       video.ID,
		URL:      *video.URL,
		Title:    *video.Title,
		Channel:  *video.Channel,
		Duration: time.Duration(*video.Duration * float64(time.Second)),
	}, nil
}
