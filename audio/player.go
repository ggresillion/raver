package audio

import (
	"fmt"
	"io"
	"log/slog"
	"time"
)

type PlaylistState int

const (
	IDLE PlaylistState = iota
	Playing
	Paused
)

func (s PlaylistState) String() string {
	switch s {
	case IDLE:
		return "IDLE"
	case Playing:
		return "PLAYING"
	case Paused:
		return "PAUSED"
	}
	return ""
}

type Player struct {
	Queue    []*Track
	State    PlaylistState
	Change   chan struct{}
	guildID  string
	autoplay bool
}

func NewPlayer(guildID string) *Player {
	return &Player{
		Queue:    make([]*Track, 0),
		Change:   make(chan struct{}),
		guildID:  guildID,
		autoplay: true,
	}
}

func (p *Player) Read(bytes []byte) (n int, err error) {
	if p.State != Playing || len(p.Queue) < 1 {
		return 0, nil
	}
	n, err = p.Queue[0].Read(bytes)
	if err == io.EOF {
		slog.Info("player: got stream end signal", "guild_id", p.guildID)
		if p.State == IDLE {
			return n, err
		}
		p.State = IDLE
		if len(p.Queue) > 1 {
			p.Queue = p.Queue[1:]
			slog.Info("player: skipped to next track", "guild_id", p.guildID)
			p.play()
		} else {
			slog.Info("player: no more tracks", "guild_id", p.guildID)
			p.Queue = make([]*Track, 0)
			p.notifyChange()
		}
		return 0, nil
	}
	return n, err
}

func (p *Player) Plug(out chan []byte) {
	go func() {
		for {
			bytes := make([]byte, 960)
			n, err := p.Read(bytes)
			if err != nil {
				slog.Info(fmt.Sprintf("[player] error writing to channel: %v", err))
				return
			}
			if n == 0 {
				time.Sleep(20 * time.Millisecond)
				continue
			}
			out <- bytes[:n]
		}
	}()
}

func (p *Player) Pause() {
	if p.State != Playing {
		return
	}
	slog.Info("player: paused", "guild_id", p.guildID)
	p.State = Paused
}

func (p *Player) Resume() {
	if p.State != Paused {
		return
	}
	slog.Info("player: resumed", "guild_id", p.guildID)
	p.State = Playing
}

func (p *Player) Skip() {
	if len(p.Queue) > 1 {
		slog.Info("player: skipping", "guild_id", p.guildID)
		t := p.Queue[0]
		p.Queue = p.Queue[1:]
		t.Close()
		p.notifyChange()
		return
	}
	slog.Info("player: cannot skip, no more track in playlist", "guild_id", p.guildID)
}

func (p *Player) Add(t *Track) error {
	p.Queue = append(p.Queue, t)
	slog.Info("player: added to queue", "track_id", t.ID, "guild_id", p.guildID)
	if p.State == IDLE {
		slog.Info("player: autoplay", "guild_id", p.guildID)
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
		slog.Info("player: cannot stop, no track in playlist", "guild_id", p.guildID)
		return
	}
	slog.Info("player: manually stopping", "guild_id", p.guildID)
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
		slog.Info("player: already playing", "guild_id", p.guildID)
		return nil
	}
	if len(p.Queue) < 1 {
		slog.Info("player: no track in playlist", "guild_id", p.guildID)
		return nil
	}
	p.State = Playing
	slog.Info("player: playing track", "track_id", p.Queue[0].ID, "guild_id", p.guildID)
	p.notifyChange()
	return nil
}

func (p *Player) notifyChange() {
	slog.Info("player: sending playlist update", "guild_id", p.guildID, "tracks", len(p.Queue), "state", p.State.String())
	go func() { p.Change <- struct{}{} }()
}
