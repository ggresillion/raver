package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ggresillion/discordsoundboard/backend/internal/config"
	"github.com/ggresillion/discordsoundboard/backend/internal/discord"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	discordAuth "github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

var (
	ErrUserNotAuthenticated = errors.New("user not authenticated")
)

const (
	AccessToken  = "access_token"
	RefreshToken = "refresh_token"
)

var conf = &oauth2.Config{
	RedirectURL:  config.Get().Host + "/api/auth/callback",
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

func GetToken(c echo.Context) string {
	return c.Get("token").(string)
}

func Authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		var token string
		cookie, err := c.Cookie(AccessToken)
		if err == nil && cookie != nil {
			token = cookie.Value
		} else {
			token = c.QueryParam(AccessToken)
		}

		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized)

		}
		c.Set("token", token)
		return next(c)
	}
}

// Login godoc
// @Summary      Login
// @Description  Redirects to login page
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "ok"
// @Failure      500  {object}  api.HTTPError
// @Router       /auth/login [get]
func (a *AuthAPI) Login(c echo.Context) error {
	uuid := uuid.New().String()
	cookie := new(http.Cookie)
	cookie.Name = "state"
	cookie.Value = uuid
	cookie.Expires = time.Now().Add(1 * time.Hour)
	c.SetCookie(cookie)

	return c.Redirect(http.StatusTemporaryRedirect, conf.AuthCodeURL(uuid))
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
	state, err := c.Cookie("state")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing state cookie")
	}

	if c.FormValue("state") != state.Value {
		return echo.NewHTTPError(http.StatusBadRequest, "")
	}

	token, err := conf.Exchange(context.Background(), c.FormValue("code"))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	accessTokenCookie := new(http.Cookie)
	accessTokenCookie.Name = AccessToken
	accessTokenCookie.Value = token.AccessToken
	accessTokenCookie.Expires = token.Expiry
	accessTokenCookie.Path = "/"
	accessTokenCookie.HttpOnly = true
	accessTokenCookie.Secure = !config.Get().Dev
	c.SetCookie(accessTokenCookie)

	refreshTokenCookie := new(http.Cookie)
	refreshTokenCookie.Name = RefreshToken
	refreshTokenCookie.Value = token.RefreshToken
	refreshTokenCookie.Path = "/"
	refreshTokenCookie.HttpOnly = true
	refreshTokenCookie.Secure = !config.Get().Dev
	c.SetCookie(refreshTokenCookie)

	return c.Redirect(http.StatusPermanentRedirect, "/")
}

// GetMe godoc
// @Summary      Get connected user
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
		var discordError *discord.ApiError
		if errors.As(err, &discordError) {
			if discordError.Code == http.StatusUnauthorized {
				return echo.NewHTTPError(http.StatusUnauthorized, discordError.Message)
			}
		}
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// Logout godoc
// @Summary      Logout
// @Description  Logout, removing the token from cookies
// @Tags         auth
// @Success      200  {object}  discord.User
// @Failure      500  {object}  api.HTTPError
// @Router       /auth/logout [get]
func (a *AuthAPI) Logout(c echo.Context) error {
	accessTokenCookie := new(http.Cookie)
	accessTokenCookie.Name = AccessToken
	accessTokenCookie.Value = ""
	accessTokenCookie.Expires = time.Now()
	accessTokenCookie.Path = "/"
	c.SetCookie(accessTokenCookie)

	refreshTokenCookie := new(http.Cookie)
	refreshTokenCookie.Name = RefreshToken
	refreshTokenCookie.Value = ""
	refreshTokenCookie.Expires = time.Now()
	refreshTokenCookie.Path = "/"
	c.SetCookie(refreshTokenCookie)

	return c.NoContent(http.StatusOK)
}
