package config

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

type Config struct {
	Dev          bool
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

	var dev bool
	if strings.HasPrefix(os.Getenv("APP_ENV"), "dev") {
		dev = true
	} else {
		dev = false
	}

	config = &Config{
		Dev:          dev,
		BotToken:     os.Getenv("BOT_TOKEN"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}
}

func Get() *Config {
	return config
}
