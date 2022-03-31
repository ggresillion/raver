package bot

import (
	"bytes"
	"encoding/json"
)

type AudioStatus int

const (
	IDLE AudioStatus = iota
	Playing
	Paused
	Buffering
	NotConnected
)

func (s AudioStatus) String() string {
	return toString[s]
}

var toString = map[AudioStatus]string{
	IDLE:         "IDLE",
	Playing:      "PLAYING",
	Paused:       "PAUSED",
	Buffering:    "BUFFERING",
	NotConnected: "NOT_CONNECTED",
}

var toID = map[string]AudioStatus{
	"IDLE":          IDLE,
	"PLAYING":       Playing,
	"PAUSED":        Paused,
	"BUFFERING":     Buffering,
	"NOT_CONNECTED": NotConnected,
}

// MarshalJSON marshals the enum as a quoted json string
func (s AudioStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *AudioStatus) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = toID[j]
	return nil
}
