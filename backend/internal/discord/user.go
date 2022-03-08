package discord

import (
	"encoding/json"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

func (dc *DiscordClient) GetUser() (*User, error) {

	res, err := dc.request("GET", "/users/@me")
	if err != nil {
		return nil, err
	}
	u := &User{}
	json.NewDecoder(res).Decode(u)
	return u, nil
}
