package web

import (
	"context"
	"io"
	"net/http"

	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

var state = "random"

type DiscordAuth struct {
	conf *oauth2.Config
}

func NewDiscordAuth() DiscordAuth {
	return DiscordAuth{conf: &oauth2.Config{
		RedirectURL: "http://localhost:3000/auth/callback",
		// This next 2 lines must be edited before running this.
		ClientID:     "id",
		ClientSecret: "secret",
		Scopes:       []string{discord.ScopeIdentify},
		Endpoint:     discord.Endpoint,
	}}
}

func (a DiscordAuth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.conf.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (a DiscordAuth) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != state {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("State does not match."))
		return
	}

	token, err := a.conf.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	res, err := a.conf.Client(context.Background(), token).Get("https://discord.com/api/users/@me")

	if err != nil || res.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte(res.Status))
		}
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(body)
}
