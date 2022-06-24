package spotify

import (
	"context"
	"errors"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/raitonoberu/ytmusic"
	"github.com/samber/lo"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/ggresillion/discordsoundboard/backend/internal/common"
	"github.com/ggresillion/discordsoundboard/backend/internal/music"
	ytdl "github.com/kkdai/youtube/v2"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyConnector struct {
	ctx    context.Context
	client *spotify.Client
}

func NewSpotifyConnector() *SpotifyConnector {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	return &SpotifyConnector{
		ctx:    ctx,
		client: client,
	}
}

// Search on spotify
func (c SpotifyConnector) Search(q string, p uint) (*music.SearchResult, error) {
	results, err := c.client.Search(c.ctx, q, spotify.SearchTypePlaylist|spotify.SearchTypeAlbum|spotify.SearchTypeArtist|spotify.SearchTypeTrack)
	if err != nil {
		return nil, err
	}

	return &music.SearchResult{
		Tracks:    getTracks(results),
		Artists:   getArtist(results),
		Albums:    getAlbums(results),
		Playlists: getPlaylists(results),
	}, nil
}

// Find track by ID
func (c SpotifyConnector) FindTrack(ID string) (*music.Track, error) {
	t, err := c.client.GetTrack(c.ctx, spotify.ID(ID))
	if err != nil {
		var sErr *spotify.Error
		if errors.As(err, sErr) {
			if sErr.Status == 404 {
				return nil, &common.NotFoundError{Err: err}
			} else {
				return nil, err
			}
		}
		return nil, err
	}

	return toTrack(t), nil
}

func (c SpotifyConnector) FindPlaylistTracks(ID string) ([]*music.Track, error) {
	playlist, err := c.client.GetPlaylistTracks(c.ctx, spotify.ID(ID))
	if err != nil {
		return nil, err
	}
	return lo.Map(playlist.Tracks, func(t spotify.PlaylistTrack, i int) *music.Track {
		return toTrack(&t.Track)
	}), nil
}

func (c SpotifyConnector) FindArtistTopTracks(ID string) ([]*music.Track, error) {
	artistTopTracks, err := c.client.GetArtistsTopTracks(c.ctx, spotify.ID(ID), "FR")
	if err != nil {
		return nil, err
	}
	return lo.Map(artistTopTracks, func(t spotify.FullTrack, i int) *music.Track {
		return toTrack(&t)
	}), nil
}

func (c SpotifyConnector) FindAlbumTracks(ID string) ([]*music.Track, error) {
	album, err := c.client.GetAlbum(c.ctx, spotify.ID(ID))
	if err != nil {
		return nil, err
	}
	return lo.Map(album.Tracks.Tracks, func(t spotify.SimpleTrack, i int) *music.Track {
		artists := make([]music.Artist, 0)
		for _, a := range t.Artists {
			artists = append(artists, music.Artist{
				ID:   string(a.ID),
				Name: a.Name,
			})
		}

		return &music.Track{
			ID:      string(t.ID),
			Title:   t.Name,
			Artists: artists,
			Album: music.Album{
				ID:        string(album.ID),
				Name:      album.Name,
				Thumbnail: getThumbnail(album.Images),
			},
			Thumbnail: getThumbnail(album.Images),
			Duration:  uint(t.Duration),
		}
	}), nil
}

// Get the stream for a given track
func (c SpotifyConnector) GetStreamURL(ID string) (string, error) {

	track, err := c.client.GetTrack(c.ctx, spotify.ID(ID))
	if err != nil {
		return "", err
	}

	// Search by track name + artists names
	q := track.Name
	for _, a := range track.Artists {
		q = q + " " + a.Name
	}

	// Get video ID from youtube
	s := ytmusic.TrackSearch(q)
	result, err := s.Next()
	if err != nil {
		return "", err
	}
	videoID := result.Tracks[0].VideoID

	// Get the video stream with ytdl
	client := ytdl.Client{}
	video, err := client.GetVideo(videoID)
	if err != nil {
		return "", err
	}

	formats := video.Formats.WithAudioChannels()

	url, err := client.GetStreamURL(video, &formats[0])
	if err != nil {
		return "", err
	}

	return url, nil
}

func getTracks(res *spotify.SearchResult) []music.Track {

	tracks := make([]music.Track, 0)
	for _, t := range res.Tracks.Tracks {

		artists := make([]music.Artist, 0)
		for _, a := range t.Artists {
			artists = append(artists, music.Artist{
				ID:   string(a.ID),
				Name: a.Name,
			})
		}

		tracks = append(tracks, music.Track{
			ID:      string(t.ID),
			Title:   t.Name,
			Artists: artists,
			Album: music.Album{
				ID:        string(t.Album.ID),
				Name:      t.Album.Name,
				Thumbnail: getThumbnail(t.Album.Images),
			},
			Thumbnail: getThumbnail(t.Album.Images),
			Duration:  uint(t.Duration),
		})
	}
	return tracks
}

func getArtist(res *spotify.SearchResult) []music.Artist {
	artists := make([]music.Artist, 0)
	for _, a := range res.Artists.Artists {
		artists = append(artists, music.Artist{
			ID:        string(a.ID),
			Name:      a.Name,
			Thumbnail: getThumbnail(a.Images),
		})
	}
	return artists
}

func getAlbums(res *spotify.SearchResult) []music.Album {

	albums := make([]music.Album, 0)
	for _, a := range res.Albums.Albums {

		artists := make([]music.Artist, 0)
		for _, a := range a.Artists {
			artists = append(artists, music.Artist{
				ID:   string(a.ID),
				Name: a.Name,
			})
		}

		albums = append(albums, music.Album{
			ID:        string(a.ID),
			Name:      a.Name,
			Thumbnail: getThumbnail(a.Images),
			Artists:   artists,
		})
	}
	return albums
}

func getPlaylists(res *spotify.SearchResult) []music.Playlist {
	playlists := make([]music.Playlist, 0)
	for _, p := range res.Playlists.Playlists {
		playlists = append(playlists, music.Playlist{
			ID:        string(p.ID),
			Name:      p.Name,
			Thumbnail: getThumbnail(p.Images),
		})
	}
	return playlists
}

func toTrack(t *spotify.FullTrack) *music.Track {
	artists := make([]music.Artist, 0)
	for _, a := range t.Artists {
		artists = append(artists, music.Artist{
			ID:   string(a.ID),
			Name: a.Name,
		})
	}

	return &music.Track{
		ID:      string(t.ID),
		Title:   t.Name,
		Artists: artists,
		Album: music.Album{
			ID:        string(t.Album.ID),
			Name:      t.Album.Name,
			Thumbnail: getThumbnail(t.Album.Images),
		},
		Thumbnail: getThumbnail(t.Album.Images),
		Duration:  uint(t.Duration),
	}
}

func getThumbnail(images []spotify.Image) string {
	if len(images) < 1 {
		return ""
	}
	return images[0].URL
}
