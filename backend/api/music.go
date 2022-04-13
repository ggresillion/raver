package api

import (
	"net/http"

	"github.com/ggresillion/discordsoundboard/backend/internal/music"
	"github.com/labstack/echo/v4"
)

type MusicAPI struct {
	manager *music.MusicPlayerManager
}

func NewMusicAPI(manager *music.MusicPlayerManager) *MusicAPI {
	return &MusicAPI{manager}
}

type AddToPlaylistPayload struct {
	ID string `json:"id"`
}

type MoveInPlaylistPayload struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type RemoveFromPlaylistPayload struct {
	Index string `json:"index"`
}

type MusicStateResponse struct {
	Playlist []music.Track `json:"playlist"`
	Status   string        `json:"status"`
}

// Search godoc
// @Summary      Search
// @Description  Searches for a song
// @Tags         music
// @Accept       json
// @Produce      json
// @Success      200  {array}   music.Track
// @Failure      400  {object}  api.HTTPError
// @Failure      404  {object}  api.HTTPError
// @Failure      500  {object}  api.HTTPError
// @Router       /music/search [get]
func (a *MusicAPI) Search(c echo.Context) error {
	q := c.QueryParam("q")

	tracks, err := a.manager.Search(q, 0)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tracks)
}

// GetState godoc
// @Summary      Get player
// @Description  Gets the music player state
// @Tags         music
// @Produce      json
// @Param        guildID  path      string                         true  "Guild ID"
// @Success      200      {object}  api.MusicStateResponse
// @Failure      400      {object}  api.HTTPError
// @Failure      404      {object}  api.HTTPError
// @Failure      500      {object}  api.HTTPError
// @Router       /guilds/{guildID}/player [get]
func (a *MusicAPI) GetState(c echo.Context) error {
	guildID := c.Param("guildID")

	player, err := a.manager.GetPlayer(guildID)
	if err != nil {
		return err
	}

	return a.returnPlayerState(c, player)
}

// AddToPlaylist godoc
// @Summary      Add to playlist
// @Description  Adds the track to the playlist
// @Tags         music
// @Accept       json
// @Produce      json
// @Param        guildID  path      string                    true  "Guild ID"
// @Param        body     body      api.AddToPlaylistPayload  true  "Body"
// @Success      200      {object}  api.MusicStateResponse
// @Failure      400      {object}  api.HTTPError
// @Failure      404      {object}  api.HTTPError
// @Failure      500      {object}  api.HTTPError
// @Router       /guilds/{guildID}/addToPlaylist [post]
func (a *MusicAPI) AddToPlaylist(c echo.Context) error {
	guildID := c.Param("guildID")

	body := &AddToPlaylistPayload{}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	player, err := a.manager.GetPlayer(guildID)
	if err != nil {
		return err
	}

	player.AddToPlaylist(body.ID)

	return a.returnPlayerState(c, player)
}

// MoveInPlaylist godoc
// @Summary      Move in playlist
// @Description  Moves a track position in playlist
// @Tags         music
// @Accept       json
// @Produce      json
// @Param        guildID  path      string                     true  "Guild ID"
// @Param        body     body      api.MoveInPlaylistPayload  true  "Body"
// @Success      200      {object}  api.MusicStateResponse
// @Failure      400      {object}  api.HTTPError
// @Failure      404      {object}  api.HTTPError
// @Failure      500      {object}  api.HTTPError
// @Router       /guilds/{guildID}/moveInPlaylist [post]
func (a *MusicAPI) MoveInPlaylist(c echo.Context) error {
	guildID := c.Param("guildID")

	body := &MoveInPlaylistPayload{}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	player, err := a.manager.GetPlayer(guildID)
	if err != nil {
		return err
	}

	player.MoveInPlaylist(body.From, body.To)

	return a.returnPlayerState(c, player)
}

// RemoveFromPlaylist godoc
// @Summary      Remove from playlist
// @Description  Removes a track from the playlist
// @Tags         music
// @Accept       json
// @Produce      json
// @Param        guildID  path      string  true  "Guild ID"
// @Param        body     body      api.RemoveFromPlaylistPayload  true  "Body"
// @Success      200      {object}  api.MusicStateResponse
// @Failure      400      {object}  api.HTTPError
// @Failure      404      {object}  api.HTTPError
// @Failure      500      {object}  api.HTTPError
// @Router       /guilds/{guildID}/removeFromPlaylist [post]
func (a *MusicAPI) RemoveFromPlaylist(c echo.Context) error {
	guildID := c.Param("guildID")

	body := &RemoveFromPlaylistPayload{}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	player, err := a.manager.GetPlayer(guildID)
	if err != nil {
		return err
	}

	player.RemoveFromPlaylist(body.Index)

	return a.returnPlayerState(c, player)
}

// Play godoc
// @Summary      Play
// @Description  Play the current track
// @Tags         music
// @Accept       json
// @Produce      json
// @Param        guildID  path      string  true  "Guild ID"
// @Success      200      {object}  api.MusicStateResponse
// @Failure      400      {object}  api.HTTPError
// @Failure      404      {object}  api.HTTPError
// @Failure      500      {object}  api.HTTPError
// @Router       /guilds/{guildID}/play [post]
func (a *MusicAPI) Play(c echo.Context) error {
	guildID := c.Param("guildID")

	player, err := a.manager.GetPlayer(guildID)
	if err != nil {
		return err
	}

	player.Play()

	return a.returnPlayerState(c, player)
}

// Pause godoc
// @Summary      Pause
// @Description  Pause the current track
// @Tags         music
// @Accept       json
// @Produce      json
// @Param        guildID  path      string  true  "Guild ID"
// @Success      200      {object}  api.MusicStateResponse
// @Failure      400      {object}  api.HTTPError
// @Failure      404      {object}  api.HTTPError
// @Failure      500      {object}  api.HTTPError
// @Router       /guilds/{guildID}/pause [post]
func (a *MusicAPI) Pause(c echo.Context) error {
	guildID := c.Param("guildID")

	player, err := a.manager.GetPlayer(guildID)
	if err != nil {
		return err
	}

	player.Play()

	return a.returnPlayerState(c, player)
}

func (a *MusicAPI) returnPlayerState(c echo.Context, player *music.MusicPlayer) error {
	response := &MusicStateResponse{
		Playlist: player.Playlist(),
		Status:   player.BotAudio().Status().String(),
	}

	return c.JSON(http.StatusOK, response)
}
