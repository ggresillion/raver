package api

import (
	"encoding/json"
	"net/http"

	"github.com/ggresillion/discordsoundboard/backend/internal/discord"
)

type DiscordAPI struct {
}

func NewDiscordAPI() *DiscordAPI {
	return &DiscordAPI{}
}

func (a *DiscordAPI) getGuilds(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	if token == nil {
		HandleUnauthorizedError(w)
		return
	}

	dc := discord.NewDiscordClient(*token)
	guilds, err := dc.GetGuilds()
	if err != nil {
		switch e := err.(type) {
		case *discord.DiscordApiError:
			w.WriteHeader(e.Code)
			w.Write([]byte(e.Message))
			return
		default:
			HandleInternalServerError(w, err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guilds)
}
