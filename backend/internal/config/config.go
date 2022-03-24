package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultYoutubeDLPath = "./youtube-dl"
)

type Config struct {
	BotToken      string
	ClientID      string
	ClientSecret  string
	YoutubeDLPath string
}

var config *Config

func init() {
	godotenv.Load()

	var youtubeDLPath string
	if os.Getenv("YOUTUBEDL_PATH") != "" {
		youtubeDLPath = os.Getenv("YOUTUBEDL_PATH")
	} else {
		youtubeDLPath = defaultYoutubeDLPath
	}

	config = &Config{
		BotToken:      os.Getenv("BOT_TOKEN"),
		ClientID:      os.Getenv("CLIENT_ID"),
		ClientSecret:  os.Getenv("CLIENT_SECRET"),
		YoutubeDLPath: youtubeDLPath,
	}
}

func Get() *Config {
	return config
}
