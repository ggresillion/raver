package api

import (
	log "github.com/sirupsen/logrus"

	"net/http"

	_ "github.com/ggresillion/discordsoundboard/backend/api/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type HTTPError struct {
	Message string `json:"message"`
}

type API struct {
	authAPI    *AuthAPI
	discordAPI *DiscordAPI
	musicAPI   *MusicAPI
	botAPI     *BotAPI
}

func NewAPI(
	authAPI *AuthAPI,
	discordAPI *DiscordAPI,
	musicAPI *MusicAPI,
	botAPI *BotAPI,
) *API {
	return &API{authAPI, discordAPI, musicAPI, botAPI}
}

// @title           DiscordSoundBoard API
// @version         1.0
// @description     This is the DiscordSoundBoard API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securitydefinitions.apikey  Authentication
// @in                          header
// @name                        Authorization
func (a *API) Listen() {
	// Create router
	e := echo.New()

	// Middlewares
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} ${status}\n",
		Output: log.StandardLogger().Writer(),
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	// Register routes
	// Public routes
	e.GET("/api/auth/login", a.authAPI.AuthLogin)
	e.GET("/api/auth/callback", a.authAPI.AuthCallback)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Static("/", "./static")

	// Authenticated routes
	r := e.Group("/api")
	r.Use(Authenticated)

	// User
	r.GET("/auth/user", a.authAPI.GetMe)

	// Guilds
	r.GET("/guilds", a.discordAPI.GetGuilds)
	r.POST("/guilds/:guildID/join", a.botAPI.JoinChannel)
	r.GET("/guilds/add", a.botAPI.AddBotToGuild)

	// Bot
	r.GET("/bot/latency", a.botAPI.GetLatency)
	r.GET("/bot/guilds", a.botAPI.GetGuilds)

	// Player
	r.GET("/guilds/:guildID/player/subscribe", a.musicAPI.SubscribeToState)
	r.GET("/guilds/:guildID/player", a.musicAPI.GetState)
	r.GET("/guilds/:guildID/player/progress", a.musicAPI.GetProgress)
	r.GET("/guilds/:guildID/player/progress/subscribe", a.musicAPI.SubscribeToProgress)
	r.POST("/guilds/:guildID/player/playlist/add", a.musicAPI.AddToPlaylist)
	r.POST("/guilds/:guildID/player/playlist/move", a.musicAPI.MoveInPlaylist)
	r.POST("/guilds/:guildID/player/playlist/remove", a.musicAPI.RemoveFromPlaylist)
	r.POST("/guilds/:guildID/player/play", a.musicAPI.Play)
	r.POST("/guilds/:guildID/player/pause", a.musicAPI.Pause)
	r.POST("/guilds/:guildID/player/stop", a.musicAPI.Stop)
	r.POST("/guilds/:guildID/player/skip", a.musicAPI.Skip)
	r.POST("/guilds/:guildID/player/time", a.musicAPI.Time)

	r.GET("/music/search", a.musicAPI.Search)

	// Start API
	log.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", e))
}
