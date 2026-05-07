package web

import (
	"context"
	"net/http"
	"raver/discord"
	"sync"
	"time"
)

type cacheEntry struct {
	user      User
	guilds    []Guild
	expires   time.Time
}

var (
	userCache   = make(map[string]cacheEntry)
	cacheMu     sync.RWMutex
	cacheTTL    = 2 * time.Minute
)

func getCachedUser(token string) (*User, bool) {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	if entry, ok := userCache[token]; ok && time.Now().Before(entry.expires) {
		return &entry.user, true
	}
	return nil, false
}

func setCachedUser(token string, user User) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	userCache[token] = cacheEntry{user: user, expires: time.Now().Add(cacheTTL)}
}

func getCachedGuilds(token string) ([]Guild, bool) {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	if entry, ok := userCache[token]; ok && time.Now().Before(entry.expires) && len(entry.guilds) > 0 {
		return entry.guilds, true
	}
	return nil, false
}

func setCachedGuilds(token string, guilds []Guild) {
	if len(guilds) == 0 {
		return
	}
	cacheMu.Lock()
	defer cacheMu.Unlock()
	if entry, ok := userCache[token]; ok {
		entry.guilds = guilds
		userCache[token] = entry
	} else {
		userCache[token] = cacheEntry{guilds: guilds, expires: time.Now().Add(cacheTTL)}
	}
}

func authenticated(bot *discord.Bot) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			var user User
			if cached, ok := getCachedUser(token.Value); ok {
				user = *cached
			} else {
				u, err := newDiscordClient(token.Value).GetUser()
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				user = *u
				setCachedUser(token.Value, user)
			}

			if user.Guild() == "" {
				var guilds []Guild
				if cached, ok := getCachedGuilds(token.Value); ok {
					guilds = cached
				} else {
					guilds, err = getGuilds(bot, token.Value)
					if err != nil {
						http.Error(w, "Unauthorized", http.StatusUnauthorized)
						return
					}
					setCachedGuilds(token.Value, guilds)
				}
				if len(guilds) > 0 {
					selectedGuilds[user.ID] = guilds[0].ID
				}
			}

			ctx := context.WithValue(r.Context(), userKey, user)
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
