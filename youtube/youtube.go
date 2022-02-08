package youtube

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

const (
	channels  int = 2     // 1 for mono, 2 for stereo
	frameSize int = 960   // uint16 size of each audio frame
	frameRate int = 48000 // audio sampling rate
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func PlayFromYoutube(v *discordgo.VoiceConnection, id string) error {

	// Send "speaking" packet over the voice websocket
	err := v.Speaking(true)
	if err != nil {
		return errors.New("couldn't set speaking")
	}

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		v.Speaking(false)
	}()

	youtube, err := getStreamFromYoutube(id)
	check(err)

	pcm, err := convertToPCM(youtube)
	check(err)

	ffmpegbuf := bufio.NewReaderSize(pcm, 16384)

	send := make(chan []int16, 2)
	defer close(send)

	close := make(chan bool)
	go func() {
		dgvoice.SendPCM(v, send)
		close <- true
	}()

	for {
		// read data from ffmpeg stdout
		audiobuf := make([]int16, frameSize*channels)
		err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		if err != nil {
			return errors.New("error reading from ffmpeg stdout")
		}

		// Send received PCM to the sendPCM channel
		select {
		case send <- audiobuf:
		case <-close:
			return nil
		}
	}
}
