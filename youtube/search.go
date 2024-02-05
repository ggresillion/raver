package youtube

import (
	"fmt"
	"log"
	"strings"
	"time"

	"raver/audio"

	"github.com/Pauloo27/searchtube"
	"github.com/buger/jsonparser"
)

func getContent(data []byte, index int) []byte {
	id := fmt.Sprintf("[%d]", index)
	contents, _, _, _ := jsonparser.Get(data, "contents", "twoColumnSearchResultsRenderer", "primaryContents", "sectionListRenderer", "contents", id, "itemSectionRenderer", "contents")
	return contents
}

func parseDurationFromMMSS(input string) (time.Duration, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid format: %s", input)
	}

	minutes := parts[0]
	seconds := parts[1]

	minutesInt := 0
	secondsInt := 0

	_, err := fmt.Sscanf(minutes, "%d", &minutesInt)
	if err != nil {
		return 0, err
	}

	_, err = fmt.Sscanf(seconds, "%d", &secondsInt)
	if err != nil {
		return 0, err
	}

	duration := time.Minute*time.Duration(minutesInt) + time.Second*time.Duration(secondsInt)
	return duration, nil
}

// Search searches for audio tracks on YouTube.
func Search(searchTerm string, limit int) (results []audio.TrackInfo, err error) {
	videos, err := searchtube.Search(searchTerm, limit)

	if err != nil {
		return nil, err
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
		results = append(results, track)
	}
	log.Printf("youtube: found %d results", len(results))
	return
}
