package goytdlp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"raver/youtube"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

type YoutubeAdapter struct {
	executable string
}

func NewYoutubeAdapter(executable string) YoutubeAdapter {
	return YoutubeAdapter{executable: executable}
}

func (y YoutubeAdapter) Name() string { return "ytdlp" }

func (y YoutubeAdapter) GetVideoByID(videoID string) (*youtube.Video, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cmd := ytdlp.New().
		SetExecutable(y.executable).
		Format("bestaudio[ext=webm][acodec=opus]").
		NoUpdate().
		NoProgress().
		NoPlaylist().
		SkipDownload().
		PrintJSON()

	slog.Info("[ytdlp] getting info for video", "video_id", videoID)
	start := time.Now()

	res, err := cmd.Run(ctx, fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID))
	if err != nil {
		return nil, err
	}

	info, err := res.GetExtractedInfo()
	if err != nil {
		return nil, err
	}
	if len(info) == 0 || info[0].URL == nil {
		return nil, errors.New("ytdlp: no stream url found")
	}

	slog.Info("[ytdlp] got stream url", "video_id", videoID, "duration", time.Since(start))

	video := info[0]

	return &youtube.Video{
		ID:       video.ID,
		URL:      *video.URL,
		Title:    deref(video.Title),
		Channel:  deref(video.Channel),
		Duration: time.Duration(derefFloat(video.Duration) * float64(time.Second)),
		IsLive:   video.IsLive != nil && *video.IsLive,
		Headers:  video.HTTPHeaders,
	}, nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
