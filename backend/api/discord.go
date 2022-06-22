package api

import (
	"net/http"

	"github.com/ggresillion/discordsoundboard/backend/internal/discord"
	echo "github.com/labstack/echo/v4"
)

type DiscordAPI struct {
}

func NewDiscordAPI() *DiscordAPI {
	return &DiscordAPI{}
}

// GetGuilds godoc
// @Summary      Get guilds
// @Description  Get the list of all guilds for the connected user
// @Tags         guilds
// @Security     Authentication
// @Accept       json
// @Produce      json
// @Success      200  {array}   discord.Guild
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /guilds [get]
func (a *DiscordAPI) GetGuilds(c echo.Context) error {
	token := getToken(c.Request())
	if token == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	dc := discord.NewDiscordClient(*token)
	guilds, err := dc.GetGuilds()
	if err != nil {
		switch e := err.(type) {
		case *discord.ApiError:
			return echo.NewHTTPError(e.Code, e.Message)
		default:
			return err
		}
	}

	return c.JSON(http.StatusOK, guilds)
}
