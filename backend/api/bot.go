package api

import (
	"net/http"

	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/discord"
	"github.com/labstack/echo/v4"
)

type LatencyResponse struct {
	// Bot latency in ms
	Latency int `json:"latency"`
}

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
// @Success      200      {string}  string
// @Failure      400      {object}  api.HTTPError
// @Failure      404      {object}  api.HTTPError
// @Failure      500      {object}  api.HTTPError
// @Router       /guilds/{guildID}/join [post]
func (a *BotAPI) JoinChannel(c echo.Context) error {
	guildID := c.Param("guildID")
	token := c.Get("token").(string)

	dc := discord.NewDiscordClient(token)

	user, err := dc.GetUser()
	if err != nil {
		return err
	}

	a.bot.GetGuildVoice(guildID).JoinUserChannel(user.ID)
	return nil
}

// GetLatency godoc
// @Summary      Get latency
// @Description  Get the server latency of the bot
// @Tags         bot
// @Security     Authentication
// @Accept       json
// @Produce      json
// @Success      200      {string}  string
// @Failure      500      {object}  api.HTTPError
// @Router       /bot/latency [post]
func (a *BotAPI) GetLatency(c echo.Context) error {
	latency := a.bot.GetLatency()
	c.JSON(http.StatusOK, &LatencyResponse{int(latency.Milliseconds())})
	return nil
}
