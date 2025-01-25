package youtube

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"raver/audio"
	"time"

	"github.com/ebml-go/webm"
	"github.com/jfbus/httprs"
)

const (
	ytDlpURL    = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
	ytDlpBinary = "yt-dlp"
)

// getVideoFromID returns a webm audio stream URL from a youtube video ID or URL
func getVideoFromID(ID string) (*Video, error) {
	ytDlp := ytDlpBinary
	// Check if yt-dlp is already installed
	if _, err := exec.LookPath(ytDlpBinary); err != nil {
		ytDlp, err = downloadYtdlp()
		if err != nil {
			return nil, err
		}
	}

	// Use the installed yt-dlp binary
	cmd := exec.Command(ytDlp, "-f", "bestaudio", "-J", ID)
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
	slog.Info("yt-dlp not found, downloading...")

	// Download yt-dlp
	resp, err := http.Get(ytDlpURL)
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

	return tmpFile.Name(), nil
}

// GetPlayableTrackFromYoutube returns a audio.Track from a given videoID
func GetPlayableTrackFromYoutube(guildID, videoID string) (*audio.Track, error) {
	v, err := getVideoFromID(videoID)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(v.URL)
	if err != nil {
		return nil, err
	}

	rs := httprs.NewHttpReadSeeker(resp)
	slog.Info("youtube: got webm video stream", "guild_id", guildID)

	audioStream, err := extractOpus(guildID, rs, resp.ContentLength)
	if err != nil {
		return nil, err
	}
	slog.Info("youtube: converting to opus stream", "guild_id", guildID)

	return audio.NewTrack(
		audio.TrackInfo{
			ID:       v.ID,
			Title:    v.Title,
			Artist:   v.Channel,
			Duration: time.Duration(v.Duration) * time.Second,
			Live:     v.IsLive,
		},
		audioStream,
	), nil
}

// extractOpus reads the incoming stream, parses it as a webm container and extract opus stream.
func extractOpus(guildID string, stream io.ReadSeeker, length int64) (*audio.AudioStream, error) {
	var w webm.WebM
	wr, err := webm.Parse(stream, &w)
	if err != nil {
		return nil, err
	}

	in := NewYTReadCloser(guildID, wr)
	slog.Info("youtube: created new input stream", "guild_id", guildID)
	audioStream := audio.NewAudioStream(guildID, in, length)

	return audioStream, nil
}

type YTReadCloser struct {
	guildID string
	wr      *webm.Reader
}

func NewYTReadCloser(guildID string, wr *webm.Reader) *YTReadCloser {
	return &YTReadCloser{guildID: guildID, wr: wr}
}

func (r *YTReadCloser) Read(bytes []byte) (n int, err error) {
	packet, ok := <-r.wr.Chan
	if !ok {
		slog.Info("youtube[%p]: closed stream", "guild_id", r.guildID)
		return 0, io.EOF
	}
	if len(packet.Data) == 0 {
		slog.Info("youtube: end of input stream", "guild_id", r.guildID)
		r.wr.Shutdown()
		return 0, io.EOF
	}
	copy(bytes, packet.Data)
	return len(packet.Data), nil
}

func (r *YTReadCloser) Close() (err error) {
	slog.Info("youtube: closing stream", "guild_id", r.guildID)
	r.wr.Shutdown()
	return
}
