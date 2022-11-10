package test

import (
	"os"
	"regexp"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/ggresillion/discordsoundboard/backend/app"
	"github.com/ggresillion/discordsoundboard/backend/internal/config"
	"github.com/joho/godotenv"
)

const projectDirName = "discordsoundboard"

func LoadEnv() {
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

func TestConnect(t *testing.T) {
	LoadEnv()

	go app.Start(&config.Config{
		BotToken:      os.Getenv("DSB_BOT_TOKEN"),
		Dev:           true,
		Host:          "localhost",
		Port:          8080,
		ClientID:      os.Getenv("DSB_CLIENT_ID"),
		ClientSecret:  os.Getenv("DSB_CLIENT_SECRET"),
		SpotifyID:     os.Getenv("DSB_SPOTIFY_ID"),
		SpotifySecret: os.Getenv("DSB_SPOTIFY_SECRET"),
	})

}
