package youtube

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"raver/audio"

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

func Search(searchTerm string, limit int) (results []audio.TrackInfo, err error) {
	url := fmt.Sprintf("https://www.youtube.com/results?search_query=%s", url.QueryEscape(searchTerm))

	var httpClient = &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Accept-Language", "en")
	res, err := httpClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	buffer, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	body := string(buffer)
	splittedScript := strings.Split(body, `window["ytInitialData"] = `)
	if len(splittedScript) != 2 {
		splittedScript = strings.Split(body, `var ytInitialData = `)
	}

	if len(splittedScript) != 2 {
		if err != nil {
			return
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
		if err != nil {
			return
		}
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
		duration, err := jsonparser.GetString(value, "videoRenderer", "lengthText", "simpleText")

		if err != nil {
			duration = ""
			live = true
		}

		parsedDuration, err := parseDurationFromMMSS(duration)
		if err != nil {
			return
		}

		results = append(results, audio.TrackInfo{
			ID:       id,
			Title:    title,
			Artist:   uploader,
			Duration: parsedDuration,
			Live:     live,
		})
	})

	return
}
