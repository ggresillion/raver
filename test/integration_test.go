package test

import (
	"raver/audio"
	"raver/youtube"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	guildID = "test"
	videoID = "VUfvRciny_Y"
)

func TestPlayerPlaySingleTrack(t *testing.T) {
	player := audio.NewPlayer(guildID)
	t.Log("getting track from youtube")
	track, err := youtube.GetPlayableTrackFromYoutube(videoID)
	if err != nil {
		t.Errorf("error loading from youtube: %v", err)
	}
	t.Log("adding track to playlist")
	player.Add(track)
	t.Log("waiting for data on player")
	bytes := make([]byte, 1000)
	n, err := player.Read(bytes)
	assert.Nil(t, err)
	assert.NotZero(t, n)
	assert.NotEmpty(t, bytes)
	assert.NotNil(t, <-player.Change)
}

func TestPlayerSkip(t *testing.T) {
	player := audio.NewPlayer(guildID)
	go func() {
		for {
			bytes := make([]byte, 1000)
			player.Read(bytes)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// get some tracks
	t.Log("adding first track")
	track1, err := youtube.GetPlayableTrackFromYoutube(videoID)
	assert.Nil(t, err)
	assert.Equal(t, audio.IDLE, player.State)
	player.Add(track1)
	t.Log("should be notified")
	<-player.Change
	assert.Equal(t, audio.Playing, player.State)
	assert.Equal(t, 1, len(player.Queue))

	t.Log("adding second track")
	track2, err := youtube.GetPlayableTrackFromYoutube(videoID)
	assert.Nil(t, err)
	player.Add(track2)
	t.Log("should be notified")
	<-player.Change
	assert.Equal(t, audio.Playing, player.State)
	assert.Equal(t, 2, len(player.Queue))
	player.Skip()
	<-player.Change
	assert.Equal(t, audio.Playing, player.State)
	assert.Equal(t, 1, len(player.Queue))
	player.Skip()
	assert.Equal(t, audio.Playing, player.State)
	assert.Equal(t, 1, len(player.Queue))
}

func TestPlayerAutostop(t *testing.T) {
	player := audio.NewPlayer(guildID)
	go func() {
		for {
			bytes := make([]byte, 1000)
			player.Read(bytes)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	track, err := youtube.GetPlayableTrackFromYoutube(videoID)
	if err != nil {
		t.Errorf("error loading from youtube: %v", err)
	}
	player.Add(track)
	<-player.Change
	assert.Equal(t, audio.Playing, player.State)
	assert.Equal(t, 1, len(player.Queue))
	<-player.Change
	assert.Equal(t, audio.IDLE, player.State)
	assert.Equal(t, 0, len(player.Queue))
}

func TestPlayerAutoskip(t *testing.T) {
	player := audio.NewPlayer(guildID)
	go func() {
		for {
			bytes := make([]byte, 1000)
			player.Read(bytes)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// get some tracks
	track1, err := youtube.GetPlayableTrackFromYoutube(videoID)
	assert.Nil(t, err)
	track2, err := youtube.GetPlayableTrackFromYoutube(videoID)
	assert.Nil(t, err)
	assert.Equal(t, audio.IDLE, player.State)
	player.Add(track1)
	<-player.Change
	assert.Equal(t, audio.Playing, player.State)
	assert.Equal(t, 1, len(player.Queue))
	player.Add(track2)
	<-player.Change
	assert.Equal(t, audio.Playing, player.State)
	assert.Equal(t, 2, len(player.Queue))
	<-player.Change
	assert.Equal(t, audio.Playing, player.State)
	assert.Equal(t, 1, len(player.Queue))
	<-player.Change
	assert.Equal(t, audio.IDLE, player.State)
	assert.Equal(t, 0, len(player.Queue))
}
