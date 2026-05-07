package web

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"raver/discord"
	"raver/youtube"
	"time"
)

//go:generate go tool templ generate

//go:embed static
var static embed.FS

func Start(bot *discord.Bot) {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(static)))

	mux.Handle("/", rootHandler(bot))
	mux.Handle("/search", authenticated(bot)(searchHandler()))
	mux.Handle("/add", authenticated(bot)(addHandler(bot)))
	mux.Handle("/player", authenticated(bot)(playerHandler(bot)))
	mux.Handle("/resume", authenticated(bot)(resumeHandler(bot)))
	mux.Handle("/pause", authenticated(bot)(pauseHandler(bot)))
	mux.Handle("/skip", authenticated(bot)(skipHandler(bot)))
	mux.Handle("/guilds", authenticated(bot)(selectGuildHandler(bot)))

	mux.HandleFunc("/auth/login", loginHandler)
	mux.HandleFunc("/auth/callback", callbackHandler)
	mux.HandleFunc("/auth/logout", logoutHandler)

	slog.Info("[web] starting server", "port", "3000")
	go http.ListenAndServe(":3000", mux)
}

type key string

const (
	userKey  key = "user"
	tokenKey key = "token"
)

func getGuilds(bot *discord.Bot, token string) ([]Guild, error) {
	client := newDiscordClient(token)
	guilds, err := client.GetUserGuilds()
	if err != nil {
		return nil, err
	}

	var guildsInCommon []Guild
	for _, g1 := range bot.Session().State.Guilds {
		for _, g2 := range guilds {
			if g1.ID == g2.ID {
				guildsInCommon = append(guildsInCommon, g2)
			}
		}
	}

	return guildsInCommon, nil
}

func rootHandler(bot *discord.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil || token.Value == "" {
			err = landing().Render(r.Context(), w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		var user User
		if cached, ok := getCachedUser(token.Value); ok {
			user = *cached
		} else {
			u, userErr := newDiscordClient(token.Value).GetUser()
			if userErr != nil {
				err = landing().Render(r.Context(), w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			user = *u
			setCachedUser(token.Value, user)
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		ctx = context.WithValue(ctx, tokenKey, token.Value)
		r = r.WithContext(ctx)

		homeHandler(bot).ServeHTTP(w, r)
	}
}

func homeHandler(bot *discord.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		var guilds []Guild
		var err error

		// Only use cache if it has guilds
		if cached, ok := getCachedGuilds(token); ok && len(cached) > 0 {
			guilds = cached
		} else {
			guilds, err = getGuilds(bot, token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if len(guilds) > 0 {
				setCachedGuilds(token, guilds)
			}
		}

		if len(guilds) == 0 {
			http.Error(w, "No guilds found. Make sure the bot is in at least one server.", http.StatusForbidden)
			return
		}

		err = index(guilds, guilds[0].ID).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func searchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.FormValue("search")

		tracks, err := youtube.Search("", query, 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, t := range tracks {
			err := track(t).Render(r.Context(), w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

func pauseHandler(bot *discord.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		guildID := r.Form.Get("guild_id")
		gbot, err := bot.Guild(guildID)
		if err != nil {
			handleError(w, err)
			return
		}

		gbot.Player.Pause()
		w.WriteHeader(http.StatusOK)
	}
}

func resumeHandler(bot *discord.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := userFromCtx(r.Context())
		gbot, err := bot.Guild(user.Guild())
		if err != nil {
			handleError(w, err)
			return
		}

		gbot.Player.Resume()
		w.WriteHeader(http.StatusOK)
	}
}

func skipHandler(bot *discord.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := userFromCtx(r.Context())
		gbot, err := bot.Guild(user.Guild())
		if err != nil {
			handleError(w, err)
			return
		}

		gbot.Player.Skip()
		w.WriteHeader(http.StatusOK)
	}
}

func playerHandler(bot *discord.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := userFromCtx(r.Context())
		slog.Info("[player] request", "user_id", user.ID, "guild", user.Guild())
		if user.Guild() == "" {
			handleError(w, errors.New("no guild selected"))
			return
		}
		slog.Info("[player] starting player stream", "guild_id", user.Guild())
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Flush support
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()

		gbot, err := bot.Guild(user.Guild())
		if err != nil {
			handleError(w, err)
			return
		}

		sendPlayerUpdate := func() {
			buf := &bytes.Buffer{}
			if err := player(gbot.Player).Render(r.Context(), buf); err != nil {
				slog.Error("render error", "err", err)
				return
			}
			slog.Info("[player] sending update", "queue_len", len(gbot.Player.Queue), "state", gbot.Player.State.String())
			fmt.Fprintf(w, "event: player\n")
			fmt.Fprintf(w, "data: %s\n\n", buf.String())
			flusher.Flush()
		}
		sendPlayerUpdate()

		for {
			select {
			case <-r.Context().Done():
				return
			case <-gbot.Player.Change:
				sendPlayerUpdate()
			case <-ticker.C:
				sendPlayerUpdate()
			}
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	slog.Error(err.Error())
}
