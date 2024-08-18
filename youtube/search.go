package youtube

import (
	"log/slog"
	"raver/audio"

	"github.com/Pauloo27/searchtube"
)

// Search searches for audio tracks on YouTube.
func Search(guildID, searchTerm string, limit int) (results []audio.TrackInfo, err error) {
	videos, err := searchtube.Search(searchTerm, limit)
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		return Search(guildID, searchTerm, limit)
	}
	for _, v := range videos {
		duration, _ := v.GetDuration()
		track := audio.TrackInfo{
			ID:       v.ID,
			Title:    v.Title,
			Artist:   v.Uploader,
			Duration: duration,
			Live:     v.Live,
		}
		if !v.Live {
			results = append(results, track)
		}
	}
	slog.Info("youtube: found tracks", "tracks", len(results), "guild_id", guildID)
	return
}
