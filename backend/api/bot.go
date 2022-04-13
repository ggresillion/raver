package api

import (
	"net/http"

	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/discord"
	"github.com/labstack/echo/v4"
)

type BotAPI struct {
	bot *bot.Bot
}

func NewBotAPI(bot *bot.Bot) *BotAPI {
	return &BotAPI{bot}
}

// JoinChannel godoc
// @Summary      Join the channel
// @Description  Join the user voice channel
// @Tags         bot
// @Security     Authentication
// @Accept       json
// @Produce      json
// @Param        guildID  path      string  true  "Guild ID"
// @Success      200      {object}  api.MusicStateResponse
// @Failure      400      {object}  api.HTTPError
// @Failure      404      {object}  api.HTTPError
// @Failure      500      {object}  api.HTTPError
// @Router       /guilds/{guildID}/join [post]
func (a *BotAPI) JoinChannel(c echo.Context) error {
	guildID := c.Param("guildID")
	token := getToken(c.Request())
	if token == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing request token")
	}

	dc := discord.NewDiscordClient(*token)

	user, err := dc.GetUser()
	if err != nil {
		return err
	}

	a.bot.GetGuildVoice(guildID).JoinUserChannel(user.ID)
	return nil
}
