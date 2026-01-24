package web

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

var state = "random"

type DiscordAuth struct {
	conf *oauth2.Config
}

func NewDiscordAuth() DiscordAuth {
	return DiscordAuth{conf: &oauth2.Config{
		RedirectURL:  "http://127.0.0.1:7331/auth/callback",
		ClientID:     os.Getenv("RAVER_CLIENT_ID"),
		ClientSecret: os.Getenv("RAVER_CLIENT_SECRET"),
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
		w.Write([]byte("state does not match."))
		return
	}

	token, err := a.conf.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type DiscordClient struct {
	Token string
}

func newDiscordClient(token string) DiscordClient {
	return DiscordClient{
		Token: token,
	}
}

type user struct {
	ID     string `json:"id"`
	Name   string `json:"username"`
	Avatar string `json:"avatar"`
}

func (c DiscordClient) GetUser(id string) (*user, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var user user
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type guild struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

func (c DiscordClient) GetUserGuilds(id string) ([]guild, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me/guilds", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var userGuilds []guild
	err = json.NewDecoder(res.Body).Decode(&userGuilds)
	if err != nil {
		return nil, err
	}

	return userGuilds, nil
}
