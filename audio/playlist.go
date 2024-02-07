package audio

import (
	"fmt"
	"log"
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
	Change  chan struct{}
	line    chan []byte
	guildID string
}

func NewPlayer(guildID string) *Player {
	return &Player{
		Queue:   make([]*Track, 0),
		Change:  make(chan struct{}),
		guildID: guildID,
		line:    make(chan []byte),
	}
}

func (p *Player) Play() error {
	if p.State == Playing {
		log.Println("player: already playing")
		return nil
	}
	if p.State == Paused {
		log.Println("player: resuming")
		p.Queue[0].Resume()
	}
	if len(p.Queue) < 1 {
		log.Println("player: no track in playlist")
		return nil
	}
	track := p.Queue[0]
	go func() {
		for {
			data, err := track.Read()
			if err != nil {
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
				return

			}
			p.line <- data
		}
	}()
	track.Play()
	p.State = Playing
	log.Printf("player: playing track %s", track.ID)
	p.notifyChange()
	return nil
}

func (p *Player) Read() []byte {
	return <-p.line
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
		p.Queue[0].Resume()
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

func (p *Player) Stop() {
	if len(p.Queue) == 0 {
		log.Println("player: cannot stop, no track in playlist")
		return
	}
	log.Println("player: manually stopping")
	t := p.Queue[0]
	p.Queue = p.Queue[:1]
	t.Stop()
}

func (p *Player) notifyChange() {
	log.Printf("player: sending playlist update (tracks: %d, state: %d)", len(p.Queue), p.State)
	go func() { p.Change <- struct{}{} }()
}
