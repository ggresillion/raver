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

func (y YoutubeConnector) Search(q string, p uint) (*music.SearchResult, error) {
	s := ytmusic.TrackSearch(q)
	result, err := s.Next()
	if err != nil {
		return nil, err
	}

	tracks := make([]music.Track, 0)
	for _, t := range result.Tracks {

		artists := make([]music.Artist, 0)
		for _, a := range t.Artists {
			artists = append(artists, music.Artist{
				ID:   a.ID,
				Name: a.Name,
			})
		}
		tracks = append(tracks, music.Track{
			ID:        string(t.VideoID),
			Title:     t.Title,
			Artists:   artists,
			Thumbnail: t.Thumbnails[0].URL,
			Album: music.Album{
				ID:   t.Album.ID,
				Name: t.Album.Name,
			},
			Duration: uint(t.Duration),
		})
	}
	return &music.SearchResult{Tracks: tracks}, nil
}

func (y YoutubeConnector) Find(ID string) (*music.Track, error) {
	s := ytmusic.TrackSearch(ID)
	result, err := s.Next()
	if err != nil {
		return nil, err
	}

	t := result.Tracks[0]

	artists := make([]music.Artist, 0)
	for _, a := range t.Artists {
		artists = append(artists, music.Artist{
			ID:   a.ID,
			Name: a.Name,
		})
	}

	return &music.Track{
		ID:        string(t.VideoID),
		Title:     t.Title,
		Artists:   artists,
		Thumbnail: t.Thumbnails[0].URL,
		Album: music.Album{
			ID:   t.Album.ID,
			Name: t.Album.Name,
		},
		Duration: uint(t.Duration),
	}, nil
}

func (y YoutubeConnector) GetStreamURL(ID string) (string, error) {

	client := ytdl.Client{}
	video, err := client.GetVideo(ID)
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
