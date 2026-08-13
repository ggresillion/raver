package test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"raver/audio"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/require"
)

// discordSilenceFrame is the fixed 3-byte opus payload Discord sends/expects for silence.
var discordSilenceFrame = []byte{0xF8, 0xFF, 0xFE}

// TestVoicePlaybackReachesDiscord joins a real voice channel with two bots: one plays a
// track through the normal pipeline, the other listens and asserts real (non-silence)
// opus packets arrive. Requires RAVER_TOKEN, RAVER_TEST_LISTENER_TOKEN, RAVER_TEST_GUILD_ID
// and RAVER_TEST_CHANNEL_ID; skipped otherwise.
func TestVoicePlaybackReachesDiscord(t *testing.T) {
	token := os.Getenv("RAVER_TOKEN")
	listenerToken := os.Getenv("RAVER_TEST_LISTENER_TOKEN")
	guild := os.Getenv("RAVER_TEST_GUILD_ID")
	channel := os.Getenv("RAVER_TEST_CHANNEL_ID")
	if token == "" || listenerToken == "" || guild == "" || channel == "" {
		t.Skip("set RAVER_TOKEN, RAVER_TEST_LISTENER_TOKEN, RAVER_TEST_GUILD_ID, RAVER_TEST_CHANNEL_ID to run this test")
	}

	player, err := newVoiceSession(token)
	require.NoError(t, err)
	defer player.Close()

	listener, err := newVoiceSession(listenerToken)
	require.NoError(t, err)
	defer listener.Close()

	playerVC, err := player.ChannelVoiceJoin(guild, channel, false, true)
	require.NoError(t, err)
	defer playerVC.Disconnect()

	listenerVC, err := listener.ChannelVoiceJoin(guild, channel, true, false)
	require.NoError(t, err)
	defer listenerVC.Disconnect()

	track, err := yt.GetPlayableTrackFromYoutube(guild, videoID)
	require.NoError(t, err)

	p := audio.NewPlayer(guild)
	p.Plug(playerVC.OpusSend)
	require.NoError(t, p.Add(track))
	defer p.Stop()

	deadline := time.After(10 * time.Second)
	for {
		select {
		case pkt := <-listenerVC.OpusRecv:
			if len(pkt.Opus) > 0 && !bytes.Equal(pkt.Opus, discordSilenceFrame) {
				return
			}
		case <-deadline:
			t.Fatal("no audio received in the voice channel within 10s")
		}
	}
}

func newVoiceSession(token string) (*discordgo.Session, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	s.Identify.Intents |= discordgo.IntentGuilds | discordgo.IntentGuildVoiceStates
	if err := s.Open(); err != nil {
		return nil, err
	}
	return s, nil
}
