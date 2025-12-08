package web

import (
	"bytes"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"raver/discord"
	"raver/youtube"
	"raver/youtube/goytdlp"

	"github.com/a-h/templ"
)

//go:generate tailwindcss -i ./app.css -o ./static/style/app.css
//go:generate go tool templ generate

//go:embed static
var static embed.FS

func Start(bot *discord.Bot) {
	auth := NewDiscordAuth()

	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(static)))

	mux.Handle("/", templ.Handler(index()))
	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/add", addHandler(bot))
	mux.HandleFunc("/player", playerHandler(bot))
	mux.HandleFunc("/resume", resumeHandler(bot))
	mux.HandleFunc("/pause", pauseHandler(bot))
	mux.HandleFunc("/skip", skipHandler(bot))

	mux.HandleFunc("/auth/login", auth.LoginHandler)
	mux.HandleFunc("/auth/callback", auth.CallbackHandler)

	slog.Info("[web] starting server", "port", "3000")
	go http.ListenAndServe(":3000", mux)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
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

func pauseHandler(bot *discord.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		gbot, err := bot.Guild(guildID)
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
		gbot, err := bot.Guild(guildID)
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
		slog.Info("player: starting player stream", "guild_id", guildID)
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

		gbot, err := bot.Guild(guildID)
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

func addHandler(bot *discord.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")

		gbot, err := bot.Guild(guildID)
		if err != nil {
			handleError(w, err)
			return
		}

		err = gbot.JoinUserChannel(userID)
		if err != nil {
			handleError(w, err)
			return
		}

		track, err := youtube.NewYoutube(goytdlp.NewYoutubeAdapter()).GetPlayableTrackFromYoutube(guildID, id)
		if err != nil {
			handleError(w, err)
			return
		}

		err = gbot.Player.Add(track)
		if err != nil {
			handleError(w, err)
			return
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	slog.Error(err.Error())
}
