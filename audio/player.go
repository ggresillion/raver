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
	Change  chan struct{}
	pause   chan bool
	line    chan []byte
	guildID string
}

func NewPlayer(guildID string) *Player {
	return &Player{
		Queue:   make([]*Track, 0),
		Change:  make(chan struct{}),
		guildID: guildID,
		line:    make(chan []byte),
		pause:   make(chan bool),
	}
}

func (p *Player) Play() error {
	if p.State != IDLE {
		log.Println("player: already playing")
		return nil
	}
	if len(p.Queue) < 1 {
		log.Println("player: no track in playlist")
		return nil
	}
	track := p.Queue[0]
	go func() {
		// track ended or stoped
		defer p.onTrackStop()
		// pipe incoming audio to the line
	play:
		log.Println("player: playing")
		p.State = Playing
		p.notifyChange()
		for {
			select {
			case paused := <-p.pause:
				if paused {
					goto pause
				}
			default:
				data, err := track.Read()
				if err != nil {
					return
				}
				p.line <- data
			}
		}
	pause:
		log.Println("player: paused")
		p.State = Paused
		p.notifyChange()
		for {
			paused := <-p.pause
			if !paused {
				goto play
			}
		}
	}()
	p.State = Playing
	log.Printf("player: playing track %s", track.ID)
	p.notifyChange()
	return nil
}

func (p *Player) Read() []byte {
	return <-p.line
}

func (p *Player) Pause() {
	if p.State != Playing {
		return
	}
	p.pause <- true
}

func (p *Player) Resume() {
	if p.State != Paused {
		return
	}
	p.pause <- false
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

func (p *Player) Progress() time.Duration {
	if len(p.Queue) < 1 {
		return 0
	}
	return p.Queue[0].Progress
}

func (p *Player) onTrackStop() {
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

}

func (p *Player) notifyChange() {
	log.Printf("player: sending playlist update (tracks: %d, state: %d)", len(p.Queue), p.State)
	go func() { p.Change <- struct{}{} }()
}
