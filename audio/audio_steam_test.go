package audio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlay(t *testing.T) {
	sampleData := []byte{255}
	in := make(chan []byte)
	s := NewAudioStream(in)
	in <- sampleData
	t.Log("play the stream")
	s.Play()
	data := <-s.Out
	t.Log("check that we received data")
	assert.Equal(t, data, sampleData)
}

func TestPause(t *testing.T) {
	sampleData := []byte{255}
	in := make(chan []byte, 1)
	s := NewAudioStream(in)
	in <- sampleData
	t.Log("play the stream")
	s.Play()
	t.Log("check that we received data")
	assert.Equal(t, sampleData, <-s.Out)
	t.Log("pause the stream")
	s.Pause()
	in <- sampleData
	t.Log("check that we received no more data")
	select {
	case <-s.Out:
		t.Error("error")
	default:
		break
	}
	t.Log("resume the steam")
	s.Pause()
	t.Log("check that we received data")
	assert.Equal(t, sampleData, <-s.Out)
}

func TestStop(t *testing.T) {
	sampleData := []byte{255}
	in := make(chan []byte)
	s := NewAudioStream(in)
	in <- sampleData
	t.Log("play the stream")
	s.Play()
	data := <-s.Out
	t.Log("check that we have data")
	assert.Equal(t, sampleData, data)
	t.Log("stop the stream")
	s.Stop()
	t.Log("check that we received stop signal")
	assert.NotNil(t, <-s.OnStop())
}
