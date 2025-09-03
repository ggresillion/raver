package web

import (
	"embed"
	"log/slog"
	"net/http"
	"raver/discord"

	"github.com/a-h/templ"
)

//go:generate tailwindcss -i ./app.css -o ./static/style/app.css
//go:generate go tool templ generate

//go:embed static
var static embed.FS

func Start(bot *discord.Bot) {
	mux := http.NewServeMux()
	mux.Handle("/", templ.Handler(index()))
	mux.Handle("/static/", http.FileServer(http.FS(static)))
	mux.HandleFunc("/search", Search)
	mux.HandleFunc("/add", Add(bot))
	slog.Info("[web] starting server", "port", "3000")
	go http.ListenAndServe(":3000", mux)
}
