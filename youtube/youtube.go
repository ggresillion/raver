package youtube

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	CHANNELS   int = 2
	FRAME_RATE int = 48000
	FRAME_SIZE int = 960
	MAX_BYTES  int = (FRAME_SIZE * 2) * 2
)

// Technically the below settings can be adjusted however that poses
// a lot of other problems that are not handled well at this time.
// These below values seem to provide the best overall performance
const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)

func getStreamFromYoutube(id string) io.Reader {
	url := "https://www.youtube.com/watch?v=" + id
	cmd := exec.Command("./youtube-dl", "-f", "251", url, "-o", "-")

	out, err := cmd.StdoutPipe()
	check(err)

	cmd.Start()

	return out
}

func convertToPCM(stream io.Reader) io.Reader {
	cmd := exec.Command("ffmpeg", "-i", "pipe:", "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	cmd.Stdin = stream
	out, err := cmd.StdoutPipe()
	check(err)

	cmd.Start()

	return out
}

func PlayFromYoutube(v *discordgo.VoiceConnection, id string) {

	stream := getStreamFromYoutube(id)

	ffmpegout := convertToPCM(stream)

	send(v, ffmpegout)
}

func PlayAudioFile(v *discordgo.VoiceConnection, stream io.Reader, stop <-chan bool) {

	ffmpegbuf := bufio.NewReaderSize(stream, 16384)

	send := make(chan []int16, 2)
	defer close(send)

	close := make(chan bool)
	go func() {
		SendPCM(v, send)
		close <- true
	}()

	for {
		// read data from ffmpeg stdout
		audiobuf := make([]int16, frameSize*channels)
		err := binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		if err != nil {
			log.Fatal("error reading from ffmpeg stdout", err)
			return
		}

		// Send received PCM to the sendPCM channel
		select {
		case send <- audiobuf:
		case <-close:
			return
		}
	}
}

func send(v *discordgo.VoiceConnection, stream io.Reader) {

	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)
	check(err)

	for {
		audiobuf := make(chan []int16, 2)
		err := binary.Read(stream, binary.LittleEndian, &audiobuf)
		check(err)
		opus, err := opusEncoder.Encode(<-audiobuf, frameSize, maxBytes)
		v.OpusSend <- opus
	}

}

// SendPCM will receive on the provied channel encode
// received PCM data into Opus then send that to Discordgo
func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		return
	}

	var err error

	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)

	if err != nil {
		log.Fatal("NewEncoder Error", err)
		return
	}

	for {

		// read pcm from chan, exit if channel is closed.
		recv, ok := <-pcm
		if !ok {
			log.Fatal("PCM Channel closed", nil)
			return
		}

		// try encoding pcm frame with Opus
		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			log.Fatal("Encoding Error", err)
			return
		}

		if v.Ready == false || v.OpusSend == nil {
			// OnError(fmt.Sprintf("Discordgo not ready for opus packets. %+v : %+v", v.Ready, v.OpusSend), nil)
			// Sending errors here might not be suited
			return
		}
		// send encoded opus data to the sendOpus channel
		v.OpusSend <- opus
	}
}

func main() {
	y := getStreamFromYoutube("kp42doFyeiM")

	io.Copy(os.Stdout, y)
}
