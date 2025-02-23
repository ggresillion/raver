package ytdlp

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"raver/youtube"
	"time"
)

var path string

//go:embed ytdlp_wrapper.py
var ytdlpWrapper embed.FS

var (
	in  io.WriteCloser
	out *bufio.Scanner
)

func init() {
	slog.Info("ytdlp: initializing...")

	// Read the embedded Python script
	scriptData, err := ytdlpWrapper.ReadFile("ytdlp_wrapper.py")
	if err != nil {
		panic(fmt.Sprintf("ytdlp: error reading embedded script: %v", err))
	}

	// Write the script to a temporary file
	tmpFile, err := os.CreateTemp("", "ytdlp_wrapper_*.py")
	if err != nil {
		panic(fmt.Sprintf("ytdlp: error creating temporary file: %v", err))
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(scriptData); err != nil {
		panic(fmt.Sprintf("ytdlp: error writing to temporary file: %v", err))
	}

	// Start the Python wrapper
	cmd := exec.Command("python3", tmpFile.Name())
	in, err = cmd.StdinPipe()
	if err != nil {
		panic(fmt.Sprintf("ytdlp: error creating stdin pipe: %v", err))
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(fmt.Sprintf("ytdlp: error creating stdout pipe: %v", err))
	}

	// Use a scanner to read from stdout
	out = bufio.NewScanner(stdout)

	// Start the process
	if err := cmd.Start(); err != nil {
		panic(fmt.Sprintf("ytdlp: error starting subprocess: %v", err))
	}

	slog.Info("ytdlp: initialized")
}

type YoutubeAdapter struct{}

func NewYoutubeAdapter() youtube.YoutubeAdapter {
	return &YoutubeAdapter{}
}

func (y *YoutubeAdapter) GetVideoByID(videoID string) (*youtube.Video, error) {
	// Send the video ID to the Python process
	_, err := fmt.Fprintln(in, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to send video ID to Python process: %w", err)
	}

	// Wait for the Python process to respond
	if !out.Scan() {
		return nil, fmt.Errorf("no response from Python process")
	}
	response := out.Text()

	// Parse the JSON response
	var video Video
	if err := json.Unmarshal([]byte(response), &video); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Map the response to the youtube.Video struct
	return &youtube.Video{
		ID:       video.ID,
		URL:      video.URL,
		Title:    video.Title,
		Channel:  video.Channel,
		Duration: time.Duration(video.Duration) * time.Second,
		IsLive:   video.IsLive,
	}, nil
}

// func downloadYtdlp() (string, error) {
// 	slog.Info("ytdlp: yt-dlp not found, downloading...")
//
// 	// Download yt-dlp
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()
//
// 	// Write the binary to a temporary file
// 	tmpFile, err := os.CreateTemp("", "yt-dlp-*")
// 	if err != nil {
// 		return "", err
// 	}
// 	defer os.Remove(tmpFile.Name()) // Clean up the temporary file
//
// 	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
// 		return "", err
// 	}
// 	if err := tmpFile.Chmod(0755); err != nil { // Make the file executable
// 		return "", err
// 	}
// 	tmpFile.Close()
//
// 	slog.Info("ytdlp: yd-dlp downloaded at " + tmpFile.Name())
// 	return tmpFile.Name(), nil
// }
