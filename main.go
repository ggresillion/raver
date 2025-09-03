package main

import (
	"log/slog"
	"os"
	"os/signal"
	"raver/discord"
	"raver/web"
	"syscall"
)

var token string

func init() {
	token = os.Getenv("BOT_TOKEN")
}

func Start(bot *discord.Bot) {
}

func main() {
	bot := discord.NewBot(token)
	err := bot.Connect()
	if err != nil {
		panic(err)
	}
	defer bot.Stop()

	web.Start(bot)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	slog.Info("exiting...")
}
