package web

import (
	"net/http"
	"raver/discord"
	"raver/youtube"
	"raver/youtube/youtubedr"
)

func addHandler(bot *discord.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := userFromCtx(r.Context())

		id := r.FormValue("id")

		gbot, err := user.GBot(bot)
		if err != nil {
			handleError(w, err)
			return
		}

		err = gbot.JoinUserChannel(user.ID)
		if err != nil {
			handleError(w, err)
			return
		}

		track, err := youtube.NewYoutube(youtubedr.NewYoutubeDRAdapter()).GetPlayableTrackFromYoutube(gbot.Guild.ID, id)
		if err != nil {
			handleError(w, err)
			return
		}

		err = gbot.Player.Add(track)
		if err != nil {
			handleError(w, err)
			return
		}
	}
}
