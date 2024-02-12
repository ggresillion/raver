package audio

import (
	"fmt"
	"io"
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
	guildID string
}

func NewPlayer(guildID string) *Player {
	return &Player{
		Queue:   make([]*Track, 0),
		Change:  make(chan struct{}),
		guildID: guildID,
	}
}

func (p *Player) Read(bytes []byte) (n int, err error) {
	if p.State != Playing || len(p.Queue) < 1 {
		return 0, nil
	}
	n, err = p.Queue[0].Read(bytes)
	if err == io.EOF {
		log.Println("player: got stream end signal")
		if p.State == IDLE {
			return
		}
		p.State = IDLE
		if len(p.Queue) > 1 {
			p.Queue = p.Queue[1:]
			log.Println("player: skipped to next track")
			p.play()
		} else {
			p.Queue = make([]*Track, 0)
			log.Println("player: no more tracks")
			p.notifyChange()
		}
		return 0, nil
	}
	return n, err
}

func (p *Player) Pause() {
	if p.State != Playing {
		return
	}
	p.State = Paused
}

func (p *Player) Resume() {
	if p.State != Paused {
		return
	}
	p.State = Playing
}
func (p *Player) Skip() {
	if len(p.Queue) > 1 {
		log.Println("player: skipping")
		t := p.Queue[0]
		p.Queue = p.Queue[1:]
		t.Close()
		return
	}
	log.Println("player: cannot skip, no more track in playlist")
}

func (p *Player) Add(t *Track) error {
	p.Queue = append(p.Queue, t)
	log.Printf("player: added %s to queue", t.ID)
	if p.State == IDLE {
		log.Println("player: autoplay")
		err := p.play()
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
	p.Queue[0].Close()
	p.Queue = []*Track{}
	p.notifyChange()
}

func (p *Player) Progress() int {
	if len(p.Queue) < 1 {
		return 0
	}
	progress := p.Queue[0].ProgressBytes
	total := p.Queue[0].TotalBytes
	if progress == 0 || total == 0 {
		return 0
	}
	return int(float64(progress) / float64(total) * 100)
}

func (p *Player) play() error {
	if p.State != IDLE {
		log.Println("player: already playing")
		return nil
	}
	if len(p.Queue) < 1 {
		log.Println("player: no track in playlist")
		return nil
	}
	p.State = Playing
	log.Printf("player: playing track %s", p.Queue[0].ID)
	p.notifyChange()
	return nil
}

func (p *Player) notifyChange() {
	log.Printf("player: sending playlist update (tracks: %d, state: %d)", len(p.Queue), p.State)
	go func() { p.Change <- struct{}{} }()
}
