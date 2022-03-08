package discord

import (
	"encoding/json"
)

type Guild struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

func (dc *DiscordClient) GetGuilds() ([]Guild, error) {
	res, err := dc.request("GET", "/users/@me/guilds")
	if err != nil {
		return nil, err
	}
	g := []Guild{}
	err = json.NewDecoder(res).Decode(&g)
	if err != nil {
		return nil, err
	}
	return g, nil
}
