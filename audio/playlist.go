package audio

import (
	"fmt"
	"log"
	"time"
)

type PlaylistState int

const (
	IDLE PlaylistState = iota
	Playing
	Paused
)

type Player struct {
	Queue   []*Track
	State   PlaylistState
	LineOut chan []byte
	Change  chan struct{}
	guildID string
}

func NewPlayer(guildID string) *Player {
	return &Player{
		Queue:   make([]*Track, 0),
		guildID: guildID,
		LineOut: make(chan []byte),
		Change:  make(chan struct{}),
	}
}

func (p *Player) Play() error {
	if p.State == Playing {
		log.Println("player: already playing")
		return nil
	}
	if p.State == Paused {
		log.Println("player: resuming")
		p.Queue[0].Play()
	}
	if len(p.Queue) < 1 {
		log.Println("player: no track in playlist")
		return nil
	}
	track := p.Queue[0]
	p.pipe(track.Out)
	track.Play()
	p.State = Playing
	log.Printf("player: playing track %s", track.ID)
	go func() {
		<-track.OnStop()
		log.Println("player: got stream end signal")
		p.State = IDLE
		p.Queue = p.Queue[1:]
		if len(p.Queue) > 0 {
			log.Println("player: skipped to next track")
			p.Play()
		} else {
			log.Println("player: no more tracks")
			p.notifyChange()
		}
	}()
	p.notifyChange()
	return nil
}

func (p *Player) Pause() {
	if len(p.Queue) < 1 {
		log.Println("player: cannot pause, no track in queue")
		return
	}
	switch p.State {
	case Playing:
		p.Queue[0].Pause()
		p.State = Paused
		log.Println("player: paused")
		p.notifyChange()
		break
	case Paused:
		p.Queue[0].Play()
		p.State = Playing
		log.Println("player: resuming")
		p.notifyChange()
	}
}

func (p *Player) Skip() {
	if len(p.Queue) > 1 {
		log.Println("player: skipping")
		p.Queue[0].Stop()
		return
	}
	log.Println("player: cannot skip, no more track in playlist")
}

func (p *Player) Add(t *Track) error {
	p.Queue = append(p.Queue, t)
	log.Printf("player: added %s to queue", t.ID)
	if p.State == IDLE {
		log.Println("player: autoplay")
		err := p.Play()
		if err != nil {
			return fmt.Errorf("player: error playing track: %v", err)
		}
	} else {
		p.notifyChange()
	}
	return nil
}

func (p *Player) pipe(c chan []byte) {
	go func() {
		for {
			data, ok := <-c
			if !ok {
				return
			}
			p.LineOut <- data
		}
	}()
}

func (p *Player) notifyChange() {
	go func() { p.Change <- struct{}{} }()
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
	*AudioStream
}

func NewTrack(infos TrackInfo, stream *AudioStream) *Track {
	return &Track{infos, stream}
}
