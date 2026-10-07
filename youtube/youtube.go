package youtube

import "time"

type Video struct {
	ID       string
	URL      string
	Title    string
	Channel  string
	Duration time.Duration
	IsLive   bool
	Headers  map[string]string
}

type YoutubeAdapter interface {
	Name() string
	GetVideoByID(videoID string) (*Video, error)
}

type Youtube struct {
	adapters []YoutubeAdapter
}

func NewYoutube(adapters ...YoutubeAdapter) Youtube {
	return Youtube{adapters: adapters}
}
