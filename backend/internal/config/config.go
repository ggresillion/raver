package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Dev           bool
	Host          string
	Port          int
	BotToken      string
	ClientID      string
	ClientSecret  string
	SpotifyID     string
	SpotifySecret string
}

var config *Config

func init() {
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

	config = &Config{
		Dev:           viper.GetBool("dev"),
		Host:          viper.GetString("host"),
		Port:          viper.GetInt("port"),
		BotToken:      viper.GetString("bot_token"),
		ClientID:      viper.GetString("client_id"),
		ClientSecret:  viper.GetString("client_secret"),
		SpotifyID:     viper.GetString("spotify_id"),
		SpotifySecret: viper.GetString("spotify_secret"),
	}
}

func Get() *Config {
	return config
}
