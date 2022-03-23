package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type API struct {
	authAPI    *AuthAPI
	discordAPI *DiscordAPI
	musicAPI   *MusicAPI
	wsAPI      *WSAPI
}

func NewAPI(
	authAPI *AuthAPI,
	discordAPI *DiscordAPI,
	musicAPI *MusicAPI,
	wsAPI *WSAPI,
) *API {
	return &API{authAPI, discordAPI, musicAPI, wsAPI}
}

func (a *API) Listen() {
	// Create router
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Register routes
	r.Get("/api/auth/login", a.authAPI.authLogin)
	r.Get("/api/auth/callback", a.authAPI.authCallback)

	r.Get("/api/guilds", a.discordAPI.getGuilds)
	r.Get("/api/guilds/{guildID}/player", a.musicAPI.getState)
	r.Post("/api/guilds/{guildID}/addToPlaylist", a.musicAPI.addToPlaylist)
	r.Post("/api/guilds/{guildID}/moveInPlaylist", a.musicAPI.moveInPlaylist)
	r.Post("/api/guilds/{guildID}/removeFromPlaylist", a.musicAPI.removeFromPlaylist)
	r.Post("/api/guilds/{guildID}/play", a.musicAPI.play)
	r.Post("/api/guilds/{guildID}/pause", a.musicAPI.pause)

	r.Get("/api/music/search", a.musicAPI.search)

	r.Get("/ws", a.wsAPI.handleConnection)

	// Start API
	log.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
