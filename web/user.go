package web

import "raver/discord"

type User struct {
	ID     string `json:"id"`
	Name   string `json:"username"`
	Avatar string `json:"avatar"`
}

func (u User) Guild() string {
	return selectedGuilds[u.ID]
}

func (u User) GBot(bot *discord.Bot) (*discord.GBot, error) {
	return bot.Guild(u.Guild())
}
