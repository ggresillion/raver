package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

var state = "random"

type DiscordAuth struct {
	conf *oauth2.Config
}

func NewDiscordAuth(origin string) DiscordAuth {
	fmt.Println(origin)
	return DiscordAuth{conf: &oauth2.Config{
		RedirectURL:  origin + "/auth/callback",
		ClientID:     os.Getenv("RAVER_CLIENT_ID"),
		ClientSecret: os.Getenv("RAVER_CLIENT_SECRET"),
		Scopes:       []string{discord.ScopeIdentify},
		Endpoint:     discord.Endpoint,
	}}
}

func getHost(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	return fmt.Sprintf("%s://%s", scheme, host)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	a := NewDiscordAuth(getHost(r))
	http.Redirect(w, r, a.conf.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	a := NewDiscordAuth(getHost(r))
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

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
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

func (c DiscordClient) GetUser() (*User, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var user User
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type Guild struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

func (c DiscordClient) GetUserGuilds() ([]Guild, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me/guilds", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	var userGuilds []Guild
	err = json.NewDecoder(res.Body).Decode(&userGuilds)
	if err != nil {
		return nil, err
	}

	return userGuilds, nil
}
