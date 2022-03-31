package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken     string
	ClientID     string
	ClientSecret string
}

var config *Config

func init() {
	godotenv.Load()

	config = &Config{
		BotToken:     os.Getenv("BOT_TOKEN"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}
}

func Get() *Config {
	return config
}
