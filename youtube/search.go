package youtube

import (
	"github.com/raitonoberu/ytmusic"
)

type Track struct {
	ID     string
	Title  string
	Artist string
}

func Search(q string) ([]Track, error) {
	s := ytmusic.Search(q)

	var tracks []Track
	for s.NextExists() {
		result, err := s.Next()
		if err != nil {
			panic(err)
		}

		for _, video := range result.Tracks {
			tracks = append(tracks, Track{
				ID:     video.VideoID,
				Title:  video.Title,
				Artist: video.Artists[0].Name,
			})
		}

	}

	return tracks, nil
}
