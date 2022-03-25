package youtube

import (
	"io"
	"os/exec"

	_ "embed"

	"github.com/ggresillion/discordsoundboard/backend/internal/config"
)

const (
	// -> opus data
	formats = "251"
)

// Returns an opus stream from a youtube video
func getStreamFromYoutube(id string) (io.Reader, error) {
	ytdlPath := config.Get().YoutubeDLPath
	cmd := exec.Command(ytdlPath, "-f", formats, id, "-o", "-")

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
