package api

import (
	"errors"
	"net/http"

	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/discord"
	"github.com/go-chi/chi/v5"
)

type BotAPI struct {
	bot *bot.Bot
}

func NewBotAPI(bot *bot.Bot) *BotAPI {
	return &BotAPI{bot}
}

func (a *BotAPI) joinChannel(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	guildID := chi.URLParam(r, "guildID")
	if token == nil {
		HandleBadRequest(w, errors.New("missing request token"))
		return
	}

	dc := discord.NewDiscordClient(*token)

	user, err := dc.GetUser()
	if err != nil {
		HandleInternalServerError(w, err)
	}

	a.bot.JoinUserChannel(guildID, user.ID)
}
