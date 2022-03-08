package youtube

import (
	"io"
	"os/exec"
)

const (
	// -> opus data
	formats  = "251"
	ytdlPath = "./youtube-dl"
)

// Returns an opus stream from a youtube video
func getStreamFromYoutube(id string) (io.Reader, error) {
	url := "https://www.youtube.com/watch?v=" + id
	cmd := exec.Command(ytdlPath, "-f", formats, url, "-o", "-")

	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return out, nil
}
