package web

import (
	"context"
	"net/http"
	"raver/discord"
)

func authenticated(bot *discord.Bot) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			user, err := newDiscordClient(token.Value).GetUser()
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if user.Guild() == "" {
				guilds, err := getGuilds(bot, token.Value)
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				selectedGuilds[user.ID] = guilds[0].ID
			}

			ctx := context.WithValue(r.Context(), userKey, *user)
			ctx = context.WithValue(ctx, tokenKey, token.Value)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func userFromCtx(ctx context.Context) User {
	u := ctx.Value(userKey)
	if u == nil {
		return User{}
	}
	return u.(User)
}

func tokenFromCtx(ctx context.Context) string {
	t := ctx.Value(tokenKey)
	if t == nil {
		return ""
	}
	return t.(string)
}
