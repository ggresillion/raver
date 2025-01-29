package youtube

import (
	"io"
	"log/slog"
	"net/http"
	"raver/audio"
	"raver/youtube/ytdlp"
	"time"

	"github.com/ebml-go/webm"
	"github.com/jfbus/httprs"
)

// GetPlayableTrackFromYoutube returns a audio.Track from a given videoID
func GetPlayableTrackFromYoutube(guildID, videoID string) (*audio.Track, error) {
	v, err := ytdlp.GetVideoByID(videoID)
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
