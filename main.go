package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"raver/discord"
	"raver/web"
	"raver/youtube"
	"raver/youtube/goytdlp"
	"raver/youtube/youtubedr"
)

var (
	token          string
	youtubeAdapter string
	ytdlpPath      string
	youtubedrProbe bool
)

func init() {
	token = os.Getenv("RAVER_TOKEN")
	youtubeAdapter = getEnv("RAVER_YOUTUBE_ADAPTERS", "youtubedr,ytdlp")
	ytdlpPath = getEnv("RAVER_YTDLP_PATH", "yt-dlp")
	youtubedrProbe = getEnv("RAVER_YOUTUBEDR_PROBE", "false") == "true"
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func youtubeAdapters(names string) ([]youtube.YoutubeAdapter, error) {
	var adapters []youtube.YoutubeAdapter
	for _, name := range strings.Split(names, ",") {
		switch strings.TrimSpace(name) {
		case "youtubedr":
			adapters = append(adapters, youtubedr.NewYoutubeDRAdapter(youtubedrProbe))
		case "ytdlp":
			adapters = append(adapters, goytdlp.NewYoutubeAdapter(ytdlpPath))
		default:
			return nil, fmt.Errorf("unknown youtube adapter %q", name)
		}
	}
	return adapters, nil
}

func main() {
	adapters, err := youtubeAdapters(youtubeAdapter)
	if err != nil {
		panic(err)
	}
	slog.Info("[youtube] adapters configured", "adapters", youtubeAdapter)

	bot := discord.NewBot(token, youtube.NewYoutube(adapters...))
	err = bot.Connect()
	if err != nil {
		panic(err)
	}
	defer bot.Stop()

	web.Start(bot)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	slog.Info("exiting...")
}
