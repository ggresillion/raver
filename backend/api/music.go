package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ggresillion/discordsoundboard/backend/internal/music"
	"github.com/gorilla/mux"
)

type MusicAPI struct {
	manager *music.MusicPlayerManager
}

func NewMusicAPI(manager *music.MusicPlayerManager) *MusicAPI {
	return &MusicAPI{manager}
}

type AddToPlaylistPayload struct {
	ID string `json:"id"`
}

func (c *MusicAPI) search(w http.ResponseWriter, r *http.Request) {
	queries, ok := r.URL.Query()["q"]

	if !ok || len(queries[0]) < 1 {
		HandleBadRequest(w, errors.New("missing \"q\" parameter"))
		return
	}

	q := queries[0]

	tracks, err := c.manager.Search(q, 0)
	if err != nil {
		HandleInternalServerError(w, err)
		return
	}

	json.NewEncoder(w).Encode(tracks)
}

func (c *MusicAPI) getState(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	guildID := params["guildID"]

	state := c.manager.GetPlayer(guildID).GetState()

	json.NewEncoder(w).Encode(state)
}

func (c *MusicAPI) addToPlaylist(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	guildID := params["guildID"]

	body := &AddToPlaylistPayload{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		HandleBadRequest(w, err)
		return
	}

	state := c.manager.GetPlayer(guildID).AddToPlaylist(body.ID)

	json.NewEncoder(w).Encode(state)
}
