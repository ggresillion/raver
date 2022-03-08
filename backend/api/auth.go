package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ggresillion/discordsoundboard/backend/internal/config"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

var state = "random"

var conf = &oauth2.Config{
	RedirectURL: "http://localhost:8080/api/auth/callback",
	// This next 2 lines must be edited before running this.
	ClientID:     config.Get().ClientID,
	ClientSecret: config.Get().ClientSecret,
	Scopes:       []string{discord.ScopeIdentify},
	Endpoint:     discord.Endpoint,
}

type AuthAPI struct {
}

func NewAuthAPI() *AuthAPI {
	return &AuthAPI{}
}

func (a *AuthAPI) authLogin(w http.ResponseWriter, r *http.Request) {
	// Step 1: Redirect to the OAuth 2.0 Authorization page.
	// This route could be named /login etc
	http.Redirect(w, r, conf.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (a *AuthAPI) authCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != state {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("State does not match."))
		return
	}
	// Step 3: We exchange the code we got for an access token
	// Then we can use the access token to do actions, limited to scopes we requested
	token, err := conf.Exchange(context.Background(), r.FormValue("code"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, fmt.Sprint("http://localhost:8080?accessToken=", token.AccessToken, "&refreshToken=", token.RefreshToken), http.StatusPermanentRedirect)
}

func getToken(r *http.Request) *string {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) < 2 {
		return nil
	}
	reqToken = splitToken[1]
	return &reqToken
}
