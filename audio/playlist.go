package audio

import (
	"log"
	"time"
)

type Playlist struct {
	Queue      []*Track
	playing    bool
	guildID    string
	onChange   func()
	bindStream func(stream *AudioStream) error
}

func NewPlaylist(guildID string, bindStream func(stream *AudioStream) error) *Playlist {
	return &Playlist{
		Queue:      make([]*Track, 0),
		guildID:    guildID,
		bindStream: bindStream,
	}
}

func (p *Playlist) Play() error {
	if p.playing {
		log.Println("Playlist: already playing")
		return nil
	}
	if len(p.Queue) < 1 {
		log.Println("Playlist: no track in playlist")
		return nil
	}
	p.playing = true
	track := p.Queue[0]
	err := p.bindStream(track.stream)
	if err != nil {
		return err
	}
	log.Printf("Playlist: playing track %s", track.ID)
	track.Play()
	go func() {
		<-p.Queue[0].stream.End
		p.playing = false
		p.Skip()
	}()
	return nil
}

func (p *Playlist) Pause() {
	if len(p.Queue) > 0 {
		p.Queue[0].stream.Pause()
		log.Println("Playlist: paused")
	}
}

func (p *Playlist) Skip() error {
	if len(p.Queue) > 1 {
		p.Queue[0].stream.Stop()
		p.Queue = p.Queue[1:]
		err := p.Play()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Playlist) Add(t *Track) error {
	p.Queue = append(p.Queue, t)
	log.Printf("Playlist: added %s to queue", t.ID)
	if !p.playing {
		log.Println("Playlist: autoplay")
		err := p.Play()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Playlist) RegisterOnChange(cb func()) {
	p.onChange = cb
}

type TrackInfo struct {
	ID       string
	Title    string
	Artist   string
	Duration time.Duration
	Live     bool
}

type Track struct {
	TrackInfo
	stream *AudioStream
}

func NewTrack(infos TrackInfo, stream *AudioStream) *Track {
	return &Track{infos, stream}
}

func (t *Track) Stream() chan []byte {
	return t.stream.Out
}

func (t *Track) Play() {
	t.stream.Play()
}

func (t *Track) Pause() {
	t.stream.Pause()
}

func (t *Track) Stop() {
	t.stream.Stop()
}
