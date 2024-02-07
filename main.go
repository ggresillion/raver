package main

import (
	"os"
	"os/signal"
	"raver/discord"
	"syscall"
)

var token string

func init() {
	token = os.Getenv("BOT_TOKEN")
}

func Start(bot *discord.Bot) {
	err := bot.Connect()
	if err != nil {
		panic(err)
	}
	defer bot.Stop()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func main() {
	bot := discord.NewBot(token)
	Start(bot)
}
