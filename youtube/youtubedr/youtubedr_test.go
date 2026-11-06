package youtubedr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamGet(t *testing.T) {
	adapter := NewYoutubeDRAdapter(true)
	v, err := adapter.GetVideoByID("HyHNuVaZJ-k")
	require.NoError(t, err)
	require.NotEmpty(t, v.URL)
}
