package main

import (
	"log"
	"os"
	"os/signal"
	"raver/discord"
	"raver/youtube"
	"syscall"

	"github.com/bwmarrin/discordgo"
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

	gb, _ := bot.Guild("423590632582414346")
	gb.JoinUserChannel("169143950203027456")
	t, _ := youtube.GetPlayableTrackFromYoutube("2ZouETXmFnM")
	gb.Player.Add(t)
	gb.PrintPlaylist(bot.Session(), &discordgo.Interaction{ChannelID: "1205166368056680498"})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("exiting...")
}

func main() {
	bot := discord.NewBot(token)
	Start(bot)
}
