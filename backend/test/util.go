package test

import (
	"os"
	"regexp"

	"github.com/ggresillion/discordsoundboard/backend/internal/config"
	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

const projectDirName = "discordsoundboard"

func loadEnv() {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.WithFields(log.Fields{
			"cause": err,
			"cwd":   cwd,
		}).Fatal("Problem loading .env file")

		os.Exit(-1)
	}
}

func TestConfig() *config.Config {
	loadEnv()
	return &config.Config{
		BotToken:      os.Getenv("DSB_BOT_TOKEN"),
		Dev:           true,
		Host:          "localhost",
		Port:          8080,
		ClientID:      os.Getenv("DSB_CLIENT_ID"),
		ClientSecret:  os.Getenv("DSB_CLIENT_SECRET"),
		SpotifyID:     os.Getenv("DSB_SPOTIFY_ID"),
		SpotifySecret: os.Getenv("DSB_SPOTIFY_SECRET"),
	}
}
