package youtube

import (
	"github.com/ggresillion/discordsoundboard/backend/internal/music"
	ytdl "github.com/kkdai/youtube/v2"
	"github.com/raitonoberu/ytmusic"
)

const (
	channels  int = 2     // 1 for mono, 2 for stereo
	frameSize int = 960   // uint16 size of each audio frame
	frameRate int = 48000 // audio sampling rate
)

type YoutubeConnector struct {
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (y YoutubeConnector) Search(q string, p uint) ([]music.Track, error) {
	s := ytmusic.TrackSearch(q)
	result, err := s.Next()
	if err != nil {
		return nil, err
	}

	tracks := make([]music.Track, 0)
	for _, t := range result.Tracks {
		var sa string
		for _, a := range t.Artists {
			sa = sa + a.Name
		}
		tracks = append(tracks, music.Track{
			ID:        string(t.VideoID),
			Title:     t.Title,
			Artist:    sa,
			Thumbnail: t.Thumbnails[0].URL,
			Album:     t.Album.Name,
			Duration:  uint(t.Duration),
			URL:       "https://www.youtube.com/watch?v=" + string(t.VideoID),
		})
	}
	return tracks, nil
}

func (y YoutubeConnector) Find(ID string) (*music.Track, error) {
	s := ytmusic.TrackSearch(ID)
	result, err := s.Next()
	if err != nil {
		return nil, err
	}

	t := result.Tracks[0]

	var sa string
	for _, a := range t.Artists {
		sa = sa + a.Name
	}
	return &music.Track{
		ID:        string(t.VideoID),
		Title:     t.Title,
		Artist:    sa,
		Thumbnail: t.Thumbnails[0].URL,
		Album:     t.Album.Name,
		Duration:  uint(t.Duration),
	}, nil
}

func (y YoutubeConnector) GetStreamURL(query string) (string, error) {

	client := ytdl.Client{}
	video, err := client.GetVideo(query)
	if err != nil {
		return "", err
	}

	formats := video.Formats.WithAudioChannels()

	url, err := client.GetStreamURL(video, &formats[0])
	if err != nil {
		return "", err
	}

	return url, nil
}
