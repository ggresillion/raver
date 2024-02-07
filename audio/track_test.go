package audio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlay(t *testing.T) {
	sampleData := []byte{255}
	s := NewAudioStream()
	s.Write(sampleData)
	t.Log("play the stream")
	s.Play()
	data, err := s.Read()
	t.Log("check that we received data")
	assert.Nil(t, err)
	assert.Equal(t, sampleData, data)
}

func TestPause(t *testing.T) {
	sampleData := []byte{255}
	s := NewAudioStream()
	s.Write(sampleData)
	t.Log("play the stream")
	s.Play()
	t.Log("check that we received data")
	data, err := s.Read()
	assert.Nil(t, err)
	assert.Equal(t, sampleData, data)
	t.Log("pause the stream")
	s.Pause()
	s.Write(sampleData)
	t.Log("check that we received no more data")
	dataCh := make(chan []byte)
	select {
	case <-dataCh:
		t.Error("should not received data when paused")
	default:
		break
	}
	t.Log("resume the steam")
	s.Play()
	t.Log("check that we received data")
	data, err = s.Read()
	assert.Nil(t, err)
	assert.Equal(t, sampleData, data)
}

func TestStop(t *testing.T) {
	sampleData := []byte{255}
	s := NewAudioStream()
	s.Write(sampleData)
	t.Log("play the stream")
	s.Play()
	t.Log("check that we have data")
	data, err := s.Read()
	assert.Nil(t, err)
	assert.Equal(t, sampleData, data)
	t.Log("stop the stream")
	s.Stop()
	t.Log("check that the stream has been stopped")
	_, err = s.Read()
	assert.Equal(t, ErrStreamClosed, err)
}
