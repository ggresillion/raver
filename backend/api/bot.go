package api

import (
	"net/http"

	"github.com/ggresillion/discordsoundboard/backend/internal/bot"
	"github.com/ggresillion/discordsoundboard/backend/internal/config"
	"github.com/ggresillion/discordsoundboard/backend/internal/discord"
	"github.com/labstack/echo/v4"
)

const permissions = "2184186880"

type LatencyResponse struct {
	// Bot latency in ms
	Latency int `json:"latency"`
}

type GuildsResponse struct {
	Guilds []discord.Guild `json:"guilds"`
}

type BotAPI struct {
	bot *bot.Bot
}

func NewBotAPI(bot *bot.Bot) *BotAPI {
	return &BotAPI{bot}
}

// GetGuilds godoc
// @Summary      Get guilds
// @Description  Get the guilds the bot and yourself are into
// @Tags         bot
// @Security     Authentication
// @Accept       json
// @Produce      json
// @Success      200      {string}  string
// @Failure      400      {object}  api.HTTPError
// @Failure      404      {object}  api.HTTPError
// @Failure      500      {object}  api.HTTPError
// @Router       /bot/guilds [post]
func (a *BotAPI) GetGuilds(c echo.Context) error {
	token := c.Get("token").(string)

	guilds, err := a.bot.GetGuilds(token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &GuildsResponse{Guilds: guilds})
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

	return a.bot.GetGuildVoice(guildID).JoinUserChannel(user.ID)
}

// AddBotToGuild godoc
// @Summary      Add bot to guild
// @Description  Add the bot to a guild
// @Tags         bot
// @Security     Authentication
// @Accept       json
// @Produce      json
// @Success      200      {string}  string
// @Failure      500      {object}  api.HTTPError
// @Router       /bot/guilds/add [post]
func (a *BotAPI) AddBotToGuild(c echo.Context) error {
	clientID := config.Get().ClientID
	return c.Redirect(http.StatusPermanentRedirect, "https://discord.com/oauth2/authorize?client_id="+clientID+"&permissions="+permissions+"&scope=bot%20applications.commands")
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
	return c.JSON(http.StatusOK, &LatencyResponse{int(latency.Milliseconds())})
}
