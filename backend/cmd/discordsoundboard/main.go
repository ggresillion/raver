package main

import (
	"log"
	"os"

	"github.com/ggresillion/discordsoundboard/backend/api"
	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/messaging"
	"github.com/ggresillion/discordsoundboard/backend/internal/music"
	"github.com/ggresillion/discordsoundboard/backend/internal/music/youtube"
)

func main() {

	// Change log output from stderr to stdout
	log.SetOutput(os.Stdout)

	hub := messaging.NewHub()

	// Start bot
	b := bot.NewBot(hub)
	err := b.StartBot()
	if err != nil {
		panic(err)
	}

	youtubeConnector := youtube.NewYoutubeConnector(hub)
	musicManager := music.NewMusicPlayerManager(youtubeConnector, hub, b)

	// Messages
	botSubscriber := bot.NewBotSubscriber(b)
	botSubscriber.SubscribeToIncomingMessages()
	musicSubscriber := music.NewMusicSubscriber(musicManager, hub)
	musicSubscriber.SubscribeToIncomingMessages()

	// Start API
	authAPI := api.NewAuthAPI()
	discordAPI := api.NewDiscordAPI()
	musicAPI := api.NewMusicAPI(musicManager)
	wsAPI := api.NewWSAPI(hub)
	botAPI := api.NewBotAPI(b)
	api := api.NewAPI(authAPI, discordAPI, musicAPI, wsAPI, botAPI)

	api.Listen()
}
