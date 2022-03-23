package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ggresillion/discordsoundboard/backend/internal/music"
	"github.com/go-chi/chi/v5"
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

type MoveInPlaylistPayload struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type RemoveFromPlaylistPayload struct {
	Index string `json:"index"`
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
	guildID := chi.URLParam(r, "guildID")

	player, err := c.manager.GetPlayer(guildID)
	if err != nil {
		HandleInternalServerError(w, err)
	}

	state := player.GetState()

	json.NewEncoder(w).Encode(state)
}

func (c *MusicAPI) addToPlaylist(w http.ResponseWriter, r *http.Request) {
	guildID := chi.URLParam(r, "guildID")

	body := &AddToPlaylistPayload{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		HandleBadRequest(w, err)
		return
	}

	player, err := c.manager.GetPlayer(guildID)
	if err != nil {
		HandleInternalServerError(w, err)
	}

	player.AddToPlaylist(body.ID)

	state := player.GetState()

	json.NewEncoder(w).Encode(state)
}

func (c *MusicAPI) moveInPlaylist(w http.ResponseWriter, r *http.Request) {
	guildID := chi.URLParam(r, "guildID")

	body := &MoveInPlaylistPayload{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		HandleBadRequest(w, err)
		return
	}

	player, err := c.manager.GetPlayer(guildID)
	if err != nil {
		HandleInternalServerError(w, err)
	}

	player.MoveInPlaylist(body.From, body.To)

	state := player.GetState()

	json.NewEncoder(w).Encode(state)
}

func (c *MusicAPI) removeFromPlaylist(w http.ResponseWriter, r *http.Request) {
	guildID := chi.URLParam(r, "guildID")

	body := &RemoveFromPlaylistPayload{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		HandleBadRequest(w, err)
		return
	}

	player, err := c.manager.GetPlayer(guildID)
	if err != nil {
		HandleInternalServerError(w, err)
	}

	player.RemoveFromPlaylist(body.Index)

	state := player.GetState()

	json.NewEncoder(w).Encode(state)
}

func (c *MusicAPI) play(w http.ResponseWriter, r *http.Request) {
	guildID := chi.URLParam(r, "guildID")

	player, err := c.manager.GetPlayer(guildID)
	if err != nil {
		HandleInternalServerError(w, err)
	}

	player.Play()

	state := player.GetState()

	json.NewEncoder(w).Encode(state)
}

func (c *MusicAPI) pause(w http.ResponseWriter, r *http.Request) {
	guildID := chi.URLParam(r, "guildID")

	player, err := c.manager.GetPlayer(guildID)
	if err != nil {
		HandleInternalServerError(w, err)
	}

	player.Play()

	state := player.GetState()

	json.NewEncoder(w).Encode(state)
}
