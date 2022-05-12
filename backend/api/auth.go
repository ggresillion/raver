package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ggresillion/discordsoundboard/backend/internal/config"
	"github.com/ggresillion/discordsoundboard/backend/internal/discord"
	"github.com/labstack/echo/v4"
	discordAuth "github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

var state = "random"

var conf = &oauth2.Config{
	RedirectURL: "http://localhost:8080/api/auth/callback",
	// This next 2 lines must be edited before running this.
	ClientID:     config.Get().ClientID,
	ClientSecret: config.Get().ClientSecret,
	Scopes:       []string{discordAuth.ScopeIdentify},
	Endpoint:     discordAuth.Endpoint,
}

type AuthAPI struct {
}

func NewAuthAPI() *AuthAPI {
	return &AuthAPI{}
}

func Authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		splitToken := strings.Split(header, "Bearer ")
		if len(splitToken) < 2 {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		token := splitToken[1]
		c.Set("token", token)
		return next(c)
	}
}

// AuthLogin godoc
// @Summary      Login
// @Description  Redirects to login page
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "ok"
// @Failure      500  {object}  api.HTTPError
// @Router       /auth/login [get]
func (a *AuthAPI) AuthLogin(c echo.Context) error {
	c.Redirect(http.StatusTemporaryRedirect, conf.AuthCodeURL(state))
	return nil
}

// AuthCallback godoc
// @Summary      Callback
// @Description  Callback that redirects to the frontend once logged in
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "ok"
// @Failure      500  {object}  api.HTTPError
// @Router       /auth/callback [get]
func (a *AuthAPI) AuthCallback(c echo.Context) error {
	if c.FormValue("state") != state {
		return echo.NewHTTPError(http.StatusBadRequest, "")
	}
	token, err := conf.Exchange(context.Background(), c.FormValue("code"))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Redirect(http.StatusPermanentRedirect, fmt.Sprint("http://localhost:3000?accessToken=", token.AccessToken, "&refreshToken=", token.RefreshToken))
	return nil
}

// GetMe godoc
// @Summary      Get connected ser
// @Description  Returns the connected user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  discord.User
// @Failure      500  {object}  api.HTTPError
// @Router       /auth/user [get]
func (a *AuthAPI) GetMe(c echo.Context) error {
	token := c.Get("token").(string)

	dc := discord.NewDiscordClient(token)

	user, err := dc.GetUser()
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, user)
	return nil
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
