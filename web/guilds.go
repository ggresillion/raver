package web

import (
	"net/http"
	"raver/discord"
)

var selectedGuilds = make(map[string]string)

func selectGuildHandler(bot *discord.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		guildID := r.FormValue("guild_id")
		user := userFromCtx(r.Context())
		selectedGuilds[user.ID] = guildID

		userGuilds, err := getGuilds(bot, tokenFromCtx(r.Context()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = guildList(userGuilds, guildID).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// func gbotForUser(ctx context.Context, bot *discord.Bot) (*discord.GBot, error) {
// 	guild, err := guildIDForUser(ctx, bot)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return bot.Guild(guild)
// }
//
// func guildIDForUser(ctx context.Context, bot *discord.Bot) (string, error) {
// 	user := userFromCtx(ctx)
// 	guild, exists := selectedGuilds[user.ID]
// 	if !exists {
// 		guilds, err := getGuilds(bot, tokenFromCtx(ctx))
// 		if err != nil {
// 			return "", err
// 		}
// 		selectedGuilds[user.ID] = guilds[0].ID
// 		guild = guilds[0].ID
// 	}
// 	return guild, nil
// }
//
// func isGuildSelected(ctx context.Context, bot *discord.Bot, guildID string) bool {
// 	selectedGuild, err := guildIDForUser(ctx, bot)
// 	if err != nil {
// 		return false
// 	}
// 	return selectedGuild == guildID
// }
