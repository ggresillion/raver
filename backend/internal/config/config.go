package config

import (
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

type Config struct {
	Dev          bool
	Host         string
	Port         int
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

	port := 8080
	if os.Getenv("PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("PORT"))
	}

	var dev bool
	if strings.HasPrefix(os.Getenv("APP_ENV"), "dev") {
		dev = true
	} else {
		dev = false
	}

	var host string
	if dev {
		host = "http://localhost:" + strconv.Itoa(port)
	} else {
		os.Getenv("HOST")
	}

	config = &Config{
		Dev:          dev,
		Host:         host,
		Port:         port,
		BotToken:     os.Getenv("BOT_TOKEN"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}
}

func Get() *Config {
	return config
}
