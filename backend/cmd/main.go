package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/ggresillion/discordsoundboard/backend/api"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/command"
	"github.com/ggresillion/discordsoundboard/backend/internal/music"
)

func init() {
	// Change log output from stderr to stdout
	log.SetOutput(os.Stdout)

	// Set the level to debug
	log.SetLevel(log.DebugLevel)
}

func main() {

	// Start bot
	b := bot.NewBot()
	err := b.StartBot()
	if err != nil {
		panic(err)
	}

	// Deps
	musicManager := music.NewPlayerManager(music.NewSpotifyConnector(), b)

	// Handle commands
	command.NewCommandHandler(b, musicManager)

	// Messages
	//botSubscriber := bot.NewBotGateway(b)
	//botSubscriber.SubscribeToIncomingMessages()
	//musicSubscriber := music.NewMusicSubscriber(musicManager, hub)
	//musicSubscriber.SubscribeToIncomingMessages()

	// Start API
	authAPI := api.NewAuthAPI()
	musicAPI := api.NewMusicAPI(musicManager)
	botAPI := api.NewBotAPI(b)
	a := api.NewAPI(authAPI, musicAPI, botAPI)

	a.Listen()
}
