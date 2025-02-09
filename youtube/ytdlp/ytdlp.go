package ytdlp

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

var (
	path       string
	cmd        *exec.Cmd
	stdinPipe  io.WriteCloser
	stdoutPipe io.ReadCloser
	mu         sync.Mutex
)

const (
	url    = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
	binary = "yt-dlp"
)

func init() {
	if _, err := exec.LookPath(binary); err != nil {
		// If not, download it
		var err error
		path, err = downloadYtdlp()
		if err != nil {
			panic(err)
		}
	} else {
		slog.Info("ytdlp: yt-dlp found")
		path = binary
	}

	// Start the persistent yt-dlp process
	startPersistentProcess()
}

func startPersistentProcess() {
	mu.Lock()
	defer mu.Unlock()

	cmd = exec.Command(path, "--no-cache-dir", "--no-progress", "-J", "-")
	var err error
	stdinPipe, err = cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdoutPipe, err = cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	slog.Info("ytdlp: started persistent yt-dlp process")
}

func GetVideoByID(ID string) (*Video, error) {
	mu.Lock()
	defer mu.Unlock()

	// Write the video ID to the stdin of the persistent process
	_, err := io.WriteString(stdinPipe, ID+"\n")
	if err != nil {
		return nil, err
	}

	// Read the JSON output from stdout
	decoder := json.NewDecoder(stdoutPipe)
	var video Video
	if err := decoder.Decode(&video); err != nil {
		return nil, err
	}

	return &video, nil
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

	slog.Info("ytdlp: yt-dlp downloaded at " + tmpFile.Name())
	return tmpFile.Name(), nil
}
