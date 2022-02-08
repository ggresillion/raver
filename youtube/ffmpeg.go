package youtube

import (
	"io"
	"os/exec"
	"strconv"
)

const (
	format = "s16le"
)

// Converts any audio stream to raw PCM data
func convertToPCM(stream io.Reader) (io.Reader, error) {
	cmd := exec.Command("ffmpeg", "-i", "pipe:", "-f", format, "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	cmd.Stdin = stream
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
