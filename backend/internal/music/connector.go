package music

import (
	"context"
	"errors"

	"github.com/raitonoberu/ytmusic"

	log "github.com/sirupsen/logrus"

	"github.com/samber/lo"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/ggresillion/discordsoundboard/backend/internal/common"
	"github.com/ggresillion/discordsoundboard/backend/internal/config"
	ytdl "github.com/kkdai/youtube/v2"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyConnector struct {
	ctx context.Context
}

var ErrCouldNotFindVideo = errors.New("could not find video on youtube")

func NewSpotifyConnector() *SpotifyConnector {
	ctx := context.Background()

	return &SpotifyConnector{
		ctx: ctx,
	}
}

func (c SpotifyConnector) GetClient() *spotify.Client {
	config := &clientcredentials.Config{
		ClientID:     config.Get().SpotifyID,
		ClientSecret: config.Get().SpotifySecret,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(c.ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(c.ctx, token)
	return spotify.New(httpClient)
}

// Search on spotify
func (c SpotifyConnector) Search(q string, p uint) (*SearchResult, error) {
	results, err := c.GetClient().Search(c.ctx, q, spotify.SearchTypePlaylist|spotify.SearchTypeAlbum|spotify.SearchTypeArtist|spotify.SearchTypeTrack)
	if err != nil {
		return nil, err
	}

	return &SearchResult{
		Tracks:    getTracks(results),
		Artists:   getArtist(results),
		Albums:    getAlbums(results),
		Playlists: getPlaylists(results),
	}, nil
}

// FindTrack by ID
func (c SpotifyConnector) FindTrack(ID string) (*Track, error) {
	t, err := c.GetClient().GetTrack(c.ctx, spotify.ID(ID))
	if err != nil {
		sErr := new(spotify.Error)
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

func (c SpotifyConnector) FindPlaylistTracks(ID string) ([]*Track, error) {
	playlist, err := c.GetClient().GetPlaylistTracks(c.ctx, spotify.ID(ID))
	if err != nil {
		return nil, err
	}
	return lo.Map(playlist.Tracks, func(t spotify.PlaylistTrack, i int) *Track {
		return toTrack(&t.Track)
	}), nil
}

func (c SpotifyConnector) FindArtistTopTracks(ID string) ([]*Track, error) {
	artistTopTracks, err := c.GetClient().GetArtistsTopTracks(c.ctx, spotify.ID(ID), "FR")
	if err != nil {
		return nil, err
	}
	return lo.Map(artistTopTracks, func(t spotify.FullTrack, i int) *Track {
		return toTrack(&t)
	}), nil
}

func (c SpotifyConnector) FindAlbumTracks(ID string) ([]*Track, error) {
	album, err := c.GetClient().GetAlbum(c.ctx, spotify.ID(ID))
	if err != nil {
		return nil, err
	}
	return lo.Map(album.Tracks.Tracks, func(t spotify.SimpleTrack, i int) *Track {
		artists := make([]Artist, 0)
		for _, a := range t.Artists {
			artists = append(artists, Artist{
				ID:   string(a.ID),
				Name: a.Name,
			})
		}

		return &Track{
			ID:      string(t.ID),
			Title:   t.Name,
			Artists: artists,
			Album: Album{
				ID:        string(album.ID),
				Name:      album.Name,
				Thumbnail: getThumbnail(album.Images),
			},
			Thumbnail: getThumbnail(album.Images),
			Duration:  uint(t.Duration),
		}
	}), nil
}

// GetStreamURL Get the stream for a given track
func (c SpotifyConnector) GetStreamURL(ID string) (string, error) {

	track, err := c.GetClient().GetTrack(c.ctx, spotify.ID(ID))
	if err != nil {
		return "", err
	}

	// Search by track name + artists names
	q := track.Name
	for _, a := range track.Artists {
		q = q + " " + a.Name
	}

	// Get video ID from YouTube
	s := ytmusic.TrackSearch(q)
	result, err := s.Next()
	if err != nil {
		return "", err
	}
	if len(result.Tracks) < 1 {
		return "", ErrCouldNotFindVideo

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

func getTracks(res *spotify.SearchResult) []Track {

	tracks := make([]Track, 0)
	for _, t := range res.Tracks.Tracks {

		artists := make([]Artist, 0)
		for _, a := range t.Artists {
			artists = append(artists, Artist{
				ID:   string(a.ID),
				Name: a.Name,
			})
		}

		tracks = append(tracks, Track{
			ID:      string(t.ID),
			Title:   t.Name,
			Artists: artists,
			Album: Album{
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

func getArtist(res *spotify.SearchResult) []Artist {
	artists := make([]Artist, 0)
	for _, a := range res.Artists.Artists {
		artists = append(artists, Artist{
			ID:        string(a.ID),
			Name:      a.Name,
			Thumbnail: getThumbnail(a.Images),
		})
	}
	return artists
}

func getAlbums(res *spotify.SearchResult) []Album {

	albums := make([]Album, 0)
	for _, a := range res.Albums.Albums {

		artists := make([]Artist, 0)
		for _, a := range a.Artists {
			artists = append(artists, Artist{
				ID:   string(a.ID),
				Name: a.Name,
			})
		}

		albums = append(albums, Album{
			ID:        string(a.ID),
			Name:      a.Name,
			Thumbnail: getThumbnail(a.Images),
			Artists:   artists,
		})
	}
	return albums
}

func getPlaylists(res *spotify.SearchResult) []Playlist {
	playlists := make([]Playlist, 0)
	for _, p := range res.Playlists.Playlists {
		playlists = append(playlists, Playlist{
			ID:        string(p.ID),
			Name:      p.Name,
			Thumbnail: getThumbnail(p.Images),
		})
	}
	return playlists
}

func toTrack(t *spotify.FullTrack) *Track {
	artists := make([]Artist, 0)
	for _, a := range t.Artists {
		artists = append(artists, Artist{
			ID:   string(a.ID),
			Name: a.Name,
		})
	}

	return &Track{
		ID:      string(t.ID),
		Title:   t.Name,
		Artists: artists,
		Album: Album{
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
