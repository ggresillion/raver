package web

import (
	"bytes"
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

	mux.Handle("/", authenticated(bot)(homeHandler(bot)))
	mux.Handle("/search", authenticated(bot)(searchHandler()))
	mux.Handle("/add", authenticated(bot)(addHandler(bot)))
	mux.Handle("/player", authenticated(bot)(playerHandler(bot)))
	mux.Handle("/resume", authenticated(bot)(resumeHandler(bot)))
	mux.Handle("/pause", authenticated(bot)(pauseHandler(bot)))
	mux.Handle("/skip", authenticated(bot)(skipHandler(bot)))
	mux.Handle("/guilds", authenticated(bot)(selectGuildHandler(bot)))

	mux.HandleFunc("/auth/login", loginHandler)
	mux.HandleFunc("/auth/callback", callbackHandler)

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

func homeHandler(bot *discord.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guilds, err := getGuilds(bot, tokenFromCtx(r.Context()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			fmt.Fprintf(w, "event: player\n")
			fmt.Fprintf(w, "data: %s\n\n", buf.String())
			flusher.Flush()
		}
		sendPlayerUpdate()

		for {
			select {
			case <-r.Context().Done():
				return
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
