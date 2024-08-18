package youtube

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"raver/audio"
	"strings"
	"time"

	"github.com/buger/jsonparser"
)

const (
	maxRetries = 3
)

// Search searches for audio tracks on YouTube.
func Search(guildID, searchTerm string, limit int) (results []audio.TrackInfo, err error) {
	for i := range maxRetries {
		results, err = performSearch(guildID, searchTerm, limit)
		if err != nil {
			slog.Error("youtube: got error on search", "error", err.Error(), "retries", i+1, "guild_id", guildID)
		}
		if err == nil && len(results) > 0 {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	slog.Info("youtube: found tracks", "tracks", len(results), "guild_id", guildID)
	return
}

func performSearch(guildID, searchTerm string, limit int) (results []audio.TrackInfo, err error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://www.youtube.com/results?search_query=%s", url.QueryEscape(searchTerm))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Cannot create GET request: %v", err)
	}
	req.Header.Add("Accept-Language", "en")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Cannot get youtube page: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	buffer, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Cannot read body: %v", err)
	}

	body := string(buffer)
	splittedScript := strings.Split(body, `window["ytInitialData"] = `)

	if len(splittedScript) != 2 {
		splittedScript = strings.Split(body, `var ytInitialData = `)
	}

	if len(splittedScript) != 2 {
		if err != nil {
			return nil, fmt.Errorf("Cannot split script: %v", err)
		}
	}

	splittedScript = strings.Split(splittedScript[1], `window["ytInitialPlayerResponse"] = null;`)
	jsonData := []byte(splittedScript[0])

	index := 0
	var contents []byte

	for {
		contents = getContent(jsonData, index)
		_, _, _, err = jsonparser.Get(contents, "[0]", "carouselAdRenderer")

		if err == nil {
			index++
		} else {
			break
		}
	}

	_, err = jsonparser.ArrayEach(contents, func(value []byte, t jsonparser.ValueType, i int, err error) {
		if limit > 0 && len(results) >= limit {
			return
		}

		id, err := jsonparser.GetString(value, "videoRenderer", "videoId")
		if err != nil {
			return
		}

		title, err := jsonparser.GetString(value, "videoRenderer", "title", "runs", "[0]", "text")
		if err != nil {
			return
		}

		uploader, err := jsonparser.GetString(value, "videoRenderer", "ownerText", "runs", "[0]", "text")
		if err != nil {
			return
		}

		live := false
		rawDuration, err := jsonparser.GetString(value, "videoRenderer", "lengthText", "simpleText")
		if err != nil {
			rawDuration = ""
			live = true
		}

		duration, _ := getDuration(rawDuration)

		results = append(results, audio.TrackInfo{
			ID:       id,
			Title:    title,
			Artist:   uploader,
			Duration: duration,
			Live:     live,
		})
	})
	return
}

func getContent(data []byte, index int) []byte {
	id := fmt.Sprintf("[%d]", index)
	contents, _, _, _ := jsonparser.Get(data, "contents", "twoColumnSearchResultsRenderer", "primaryContents", "sectionListRenderer", "contents", id, "itemSectionRenderer", "contents")
	return contents
}

func getDuration(rawDuration string) (duration time.Duration, err error) {
	if duration == 0 {
		str := rawDuration + "s"
		if strings.Count(str, ":") == 2 {
			str = strings.Replace(str, ":", "h", 1)
		}
		str = strings.Replace(str, ":", "m", 1)

		duration, err = time.ParseDuration(str)
	}
	return
}
