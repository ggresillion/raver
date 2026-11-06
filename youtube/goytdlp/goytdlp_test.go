package goytdlp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamGet(t *testing.T) {
	adapter := NewYoutubeAdapter("")
	v, err := adapter.GetVideoByID("HyHNuVaZJ-k")
	require.NoError(t, err)
	require.NotEmpty(t, v.URL)
}
