package music

import (
	"bytes"
	"encoding/json"
)

type MusicPlayerStatus int

const (
	IDLE MusicPlayerStatus = iota
	Playing
	Paused
	Buffering
)

func (s MusicPlayerStatus) String() string {
	return toString[s]
}

var toString = map[MusicPlayerStatus]string{
	IDLE:      "IDLE",
	Playing:   "PLAYING",
	Paused:    "PAUSED",
	Buffering: "BUFFERING",
}

var toID = map[string]MusicPlayerStatus{
	"IDLE":      IDLE,
	"PLAYING":   Playing,
	"PAUSED":    Paused,
	"BUFFERING": Buffering,
}

// MarshalJSON marshals the enum as a quoted json string
func (s MusicPlayerStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *MusicPlayerStatus) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = toID[j]
	return nil
}
