package ytdlp

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
)

var path string

const (
	url    = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
	binary = "yt-dlp"
)

func init() {
	if _, err := exec.LookPath(binary); err != nil {
		// If not, download it
		path, err = downloadYtdlp()
		if err != nil {
			panic(err)
		}
		return
	}
	slog.Info("ytdlp: yt-dlp found")
	path = binary
}

func GetVideoByID(ID string) (*Video, error) {
	// Use the installed yt-dlp binary
	cmd := exec.Command(path, "-f", "bestaudio", "-J", ID)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var video *Video
	err = json.Unmarshal(output, &video)
	if err != nil {
		return nil, err
	}

	return video, nil
}

func downloadYtdlp() (string, error) {
	slog.Info("ytdlp: yt-dlp not found, downloading...")

	// Download yt-dlp
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Write the binary to a temporary file
	tmpFile, err := os.CreateTemp("", "yt-dlp-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name()) // Clean up the temporary file

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return "", err
	}
	if err := tmpFile.Chmod(0755); err != nil { // Make the file executable
		return "", err
	}
	tmpFile.Close()

	slog.Info("ytdlp: yd-dlp downloaded at " + tmpFile.Name())
	return tmpFile.Name(), nil
}
