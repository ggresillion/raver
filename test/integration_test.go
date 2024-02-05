package test

import (
	"raver/audio"
	"raver/youtube"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	guildID = "test"
	videoID = "hqS4u4aiZ8g"
)

func TestPlayerPlaySingleTrack(t *testing.T) {
	player := audio.NewPlayer(guildID)
	change := player.OnChange()
	track, err := youtube.GetPlayableTrackFromYoutube(videoID)
	if err != nil {
		t.Errorf("error loading from youtube: %v", err)
	}
	player.Add(track)
	<-player.LineOut
	assert.NotNil(t, <-change)
}

func TestPlayerSkip(t *testing.T) {
	player := audio.NewPlayer(guildID)
	change := player.OnChange()
	go func() {
		for {
			<-player.LineOut
		}
	}()

	// get some tracks
	track1, err := youtube.GetPlayableTrackFromYoutube(videoID)
	assert.Nil(t, err)
	assert.Equal(t, audio.IDLE, player.State)
	player.Add(track1)
	<-change
	assert.Equal(t, audio.Playing, player.State)
	<-change

	track2, err := youtube.GetPlayableTrackFromYoutube(videoID)
	assert.Nil(t, err)
	player.Add(track2)
	assert.Equal(t, 2, len(player.Queue))
	player.Skip()
	<-change
	<-change
	assert.Equal(t, audio.Playing, player.State)
	assert.Equal(t, 1, len(player.Queue))
}
