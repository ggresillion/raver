package config

import (
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host         string
	BotToken     string
	ClientID     string
	ClientSecret string
}

var config *Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Debugf("error loading config %e", err)
	}

	var host string
	if os.Getenv("HOST") == "" {
		host = "http://localhost:8080"
	} else {
		host = "https://" + os.Getenv("HOST")
	}

	config = &Config{
		Host:         host,
		BotToken:     os.Getenv("BOT_TOKEN"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}
}

func Get() *Config {
	return config
}
