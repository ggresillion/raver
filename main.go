package main

import (
	"os"
	"os/signal"
	"raver/api"
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
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		bot.Stop()
		os.Exit(0)
	}()
	err = api.Listen()
	if err != nil {
		panic(err)
	}
}

func main() {
	bot := discord.NewBot(token)
	Start(bot)
}
