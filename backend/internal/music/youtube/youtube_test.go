package youtube

import (
	"io"
	"os"
	"testing"
)

func copyToFile(r io.Reader, filename string) {

	f, err := os.Create(filename)
	check(err)

	_, err = io.Copy(f, r)
	check(err)
}

func readFromFile(filename string) io.Reader {
	f, err := os.Open(filename)
	check(err)
	return f
}

func TestYoutubeDownload(t *testing.T) {

	s, err := getStreamFromYoutube("tUlthCngK9U")
	check(err)

	copyToFile(s, "output.opus")
}

func TestConvertToPCM(t *testing.T) {

	in := readFromFile("output.opus")

	r, err := convertToPCM(in)
	check(err)

	copyToFile(r, "output.pcm")
}
