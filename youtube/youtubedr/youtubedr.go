package youtubedr

import (
	"fmt"
	"log/slog"
	"net/http"
	"raver/youtube"
	"strings"

	ytdr "github.com/kkdai/youtube/v2"
)

type YoutubeDR struct {
	probe bool
}

func init() {
	ytdr.DefaultClient = ytdr.ClientInfo{
		Name:        "VISIONOS",
		Version:     "1.02",
		DeviceModel: "RealityDevice17,1",
		UserAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 15_7_3) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.0 Safari/605.1.15",
	}
}

func (y *YoutubeDR) Name() string { return "youtubedr" }

func NewYoutubeDRAdapter(probe bool) youtube.YoutubeAdapter {
	return &YoutubeDR{probe: probe}
}

func (y *YoutubeDR) GetVideoByID(videoID string) (*youtube.Video, error) {
	yt := ytdr.Client{}
	video, err := yt.GetVideo(videoID)
	if err != nil {
		return nil, err
	}

	var url string
	var format *ytdr.Format
	for _, f := range video.Formats.WithAudioChannels() {
		if !strings.Contains(f.MimeType, "opus") {
			continue
		}
		u, err := yt.GetStreamURL(video, &f)
		if err != nil {
			continue
		}
		if y.probe && !isStreamable(u, f.ContentLength) {
			slog.Info("[youtubedr] stream url not available", "video_id", videoID, "itag", f.ItagNo, "bitrate", f.Bitrate)
			continue
		}
		url, format = u, &f
		break
	}

	if format == nil {
		return nil, fmt.Errorf("%w: no streamable opus format found", youtube.ErrStreamRefused)
	}

	slog.Info("[youtubedr] got stream url", "video_id", videoID, "client", ytdr.DefaultClient.Name, "itag", format.ItagNo, "format", format.MimeType)

	return &youtube.Video{
		ID:       video.ID,
		URL:      url,
		Title:    video.Title,
		Channel:  video.Author,
		Duration: video.Duration,
		IsLive:   false,
	}, nil
}

func isStreamable(url string, length int64) bool {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", length-1, length-1))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusPartialContent
}
