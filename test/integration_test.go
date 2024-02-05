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
	track, err := youtube.GetPlayableTrackFromYoutube(videoID)
	if err != nil {
		t.Errorf("error loading from youtube: %v", err)
	}
	player.Add(track)
	<-player.LineOut
	assert.NotNil(t, <-player.Change)
}

func TestPlayerSkip(t *testing.T) {
	player := audio.NewPlayer(guildID)
	go func() {
		for {
			<-player.LineOut
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// get some tracks
	track1, err := youtube.GetPlayableTrackFromYoutube(videoID)
	assert.Nil(t, err)
	assert.Equal(t, audio.IDLE, player.State)
	player.Add(track1)
	<-player.Change
	assert.Equal(t, audio.Playing, player.State)
	assert.Equal(t, 1, len(player.Queue))

	track2, err := youtube.GetPlayableTrackFromYoutube(videoID)
	assert.Nil(t, err)
	player.Add(track2)
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
			<-player.LineOut
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
			<-player.LineOut
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
	player.Add(track2)
	<-player.Change
	<-player.Change
	assert.Equal(t, audio.Playing, player.State)
	assert.Equal(t, 2, len(player.Queue))
	<-player.Change
	<-player.Change
	assert.Equal(t, audio.Playing, player.State)
	assert.Equal(t, 1, len(player.Queue))
	<-player.Change
	assert.Equal(t, audio.IDLE, player.State)
	assert.Equal(t, 0, len(player.Queue))
}
