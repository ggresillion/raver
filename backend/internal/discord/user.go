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

func (c *Client) GetUser() (*User, error) {

	res, err := c.request("GET", "/users/@me")
	if err != nil {
		return nil, err
	}
	u := &User{}
	err = json.NewDecoder(res).Decode(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
