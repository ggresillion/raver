package youtube

import "time"

type Video struct {
	ID       string
	URL      string
	Title    string
	Channel  string
	Duration time.Duration
	IsLive   bool
}

type YoutubeAdapter interface {
	GetVideoByID(videoID string) (*Video, error)
}

type Youtube struct {
	adapter YoutubeAdapter
}

func NewYoutube(adapter YoutubeAdapter) Youtube {
	return Youtube{adapter: adapter}
}
