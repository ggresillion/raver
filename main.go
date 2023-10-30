package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"raver/discord"
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
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-stop
	log.Println("Gracefully shutting down")
	bot.Stop()
}

func main() {
	bot := discord.NewBot(token)
	Start(bot)
}
