package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"raver/discord"
)

var (
	token string
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	bot := discord.NewBot(token)
	err := bot.Connect()
	if err != nil {
		panic(err)
	}

	// tracks, err := youtube.Search("Daft punk")
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(tracks)

	// audioStream, err := youtube.GetStreamFromYoutube(tracks[1].ID)
	// if err != nil {
	// 	panic(err)
	// }

	// bot.PlayStream(audioStream, "423590632582414346", "423590632582414350")

	// audioStream.BlockUntilStop()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-stop

	log.Println("Gracefully shutting down")
	bot.Stop()

}
