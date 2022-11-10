package app

import (
	"github.com/ggresillion/discordsoundboard/backend/api"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/command"
	"github.com/ggresillion/discordsoundboard/backend/internal/config"
	"github.com/ggresillion/discordsoundboard/backend/internal/music"
)

func Start(conf *config.Config) {
	// Start bot
	b := bot.NewBot(conf.BotToken)
	err := b.StartBot()
	if err != nil {
		panic(err)
	}

	// Deps
	musicManager := music.NewPlayerManager(music.NewSpotifyConnector(conf.SpotifyID, conf.SpotifySecret), b)

	// Handle commands
	command.NewCommandHandler(b, musicManager)

	// Start API
	authAPI := api.NewAuthAPI(conf.Host, conf.ClientID, conf.ClientSecret, conf.Dev)
	musicAPI := api.NewMusicAPI(musicManager)
	botAPI := api.NewBotAPI(conf.ClientID, b)
	a := api.NewAPI(conf.Dev, conf.Port, authAPI, musicAPI, botAPI)

	a.Listen()
}
