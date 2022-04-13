package api

import (
	"log"
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
	wsAPI      *WSAPI
	botAPI     *BotAPI
}

func NewAPI(
	authAPI *AuthAPI,
	discordAPI *DiscordAPI,
	musicAPI *MusicAPI,
	wsAPI *WSAPI,
	botAPI *BotAPI,
) *API {
	return &API{authAPI, discordAPI, musicAPI, wsAPI, botAPI}
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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://labstack.com", "https://labstack.net"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Register routes
	e.GET("/api/auth/login", a.authAPI.AuthLogin)
	e.GET("/api/auth/callback", a.authAPI.AuthCallback)

	e.GET("/api/guilds", a.discordAPI.GetGuilds)
	e.GET("/api/guilds/:guildID/player", a.musicAPI.GetState)
	e.POST("/api/guilds/:guildID/addToPlaylist", a.musicAPI.AddToPlaylist)
	e.POST("/api/guilds/:guildID/moveInPlaylist", a.musicAPI.MoveInPlaylist)
	e.POST("/api/guilds/:guildID/removeFromPlaylist", a.musicAPI.RemoveFromPlaylist)
	e.POST("/api/guilds/:guildID/play", a.musicAPI.Play)
	e.POST("/api/guilds/:guildID/pause", a.musicAPI.Pause)
	e.POST("/api/guilds/:guildID/join", a.botAPI.JoinChannel)

	e.GET("/api/music/search", a.musicAPI.Search)

	// e.GET("/ws", a.wsAPI.handleConnection)

	e.Static("/", "./static")

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start API
	log.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", e))
}
