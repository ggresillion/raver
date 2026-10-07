package test

import (
	"io"
	"os/exec"
	"raver/youtube"
	"raver/youtube/goytdlp"
	"raver/youtube/youtubedr"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var adapterVideos = map[string]string{
	"official music video": "kJQP7kiw5Fk",
	"regular video":        "aqz-KE-bpKQ",
}

func youtubeAdapters(tb testing.TB) []youtube.YoutubeAdapter {
	if testing.Short() {
		tb.Skip("requires network access to youtube")
	}
	adapters := []youtube.YoutubeAdapter{youtubedr.NewYoutubeDRAdapter(false)}
	if _, err := exec.LookPath(ytdlpPath()); err == nil {
		adapters = append(adapters, goytdlp.NewYoutubeAdapter(ytdlpPath()))
	} else {
		tb.Logf("yt-dlp not found at %q, skipping ytdlp adapter, set RAVER_YTDLP_PATH", ytdlpPath())
	}
	return adapters
}

func TestYoutubeAdapters(t *testing.T) {
	for _, adapter := range youtubeAdapters(t) {
		t.Run(adapter.Name(), func(t *testing.T) {
			for name, id := range adapterVideos {
				t.Run(name, func(t *testing.T) {
					track, err := youtube.NewYoutube(adapter).GetPlayableTrackFromYoutube(guildID, id)
					require.NoError(t, err)
					assert.NotEmpty(t, track.Title)
					assert.NotZero(t, track.Duration)

					n, err := io.Copy(io.Discard, track.AudioStream)
					require.NoError(t, err)
					assert.Greater(t, n, int64(1<<20))
				})
			}
		})
	}
}

func BenchmarkYoutubeAdapters(b *testing.B) {
	for _, adapter := range youtubeAdapters(b) {
		for name, id := range adapterVideos {
			b.Run(adapter.Name()+"/"+name, func(b *testing.B) {
				yt := youtube.NewYoutube(adapter)
				fails := 0
				for b.Loop() {
					track, err := yt.GetPlayableTrackFromYoutube(guildID, id)
					if err != nil {
						b.Log(err)
						fails++
						continue
					}
					track.AudioStream.Close()
				}
				b.ReportMetric(float64(fails)/float64(b.N), "fails/op")
			})
		}
	}
}
