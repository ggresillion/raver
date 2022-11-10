package main

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ggresillion/discordsoundboard/backend/app"
	"github.com/ggresillion/discordsoundboard/backend/internal/config"
)

func init() {
	// Change log output from stderr to stdout
	log.SetOutput(os.Stdout)

	// Set the level to debug
	log.SetLevel(log.DebugLevel)

	godotenv.Load()

	viper.AddConfigPath(".")

	viper.SetEnvPrefix("dsb")
	viper.BindEnv("dev")
	viper.BindEnv("host")
	viper.BindEnv("port")
	viper.BindEnv("bot_token")
	viper.BindEnv("client_id")
	viper.BindEnv("client_secret")
	viper.BindEnv("spotify_id")
	viper.BindEnv("spotify_secret")

	viper.ReadInConfig()

	viper.SetDefault("dev", true)
	viper.SetDefault("host", "http://localhost:8080")
	viper.SetDefault("port", 8080)
}

func main() {
	app.Start(&config.Config{
		Dev:           viper.GetBool("dev"),
		Host:          viper.GetString("host"),
		Port:          viper.GetInt("port"),
		BotToken:      viper.GetString("bot_token"),
		ClientID:      viper.GetString("client_id"),
		ClientSecret:  viper.GetString("client_secret"),
		SpotifyID:     viper.GetString("spotify_id"),
		SpotifySecret: viper.GetString("spotify_secret"),
	})
}
