package spotify

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Client struct {
	spotify *spotify.Client
	logger  *zap.SugaredLogger
}

func New(token *oauth2.Token) Client {
	return Client{
		spotify: newSpotifyClient(context.Background(), token),
		logger:  zap.L().Sugar(),
	}
}

func (c Client) CreatePlaylist(ctx context.Context, name string, artists []string) (spotify.ID, error) {
	playlistID, err := c.createPlaylist(ctx, name)
	if err != nil {
		return "", fmt.Errorf("creating playlist: %w", err)
	}
	err = c.addNewTracksToPlaylist(ctx, playlistID, artists)
	if err != nil {
		return "", err
	}
	return playlistID, nil
}

func (c Client) AppendToPlaylist(ctx context.Context, playlistID string, artists []string) error {
	return c.addNewTracksToPlaylist(ctx, spotify.ID(playlistID), artists)
}

func (c Client) addNewTracksToPlaylist(ctx context.Context, playlistID spotify.ID, artists []string) error {
	var newTracks []spotify.SimpleTrack
	for _, artist := range artists {
		tracks, err := c.getNewTracksOfArtist(ctx, artist)
		if err != nil {
			return fmt.Errorf("getting new tracks of artist %v: %w", artist, err)
		}
		c.logger.Debugf("found %d tracks of artist %v", len(tracks), artist)
		newTracks = append(newTracks, tracks...)
	}
	err := c.addTracksToPlaylist(ctx, playlistID, newTracks)
	if err != nil {
		return fmt.Errorf("adding tracks to playlist: %w", err)
	}
	return nil
}

func (c Client) createPlaylist(ctx context.Context, name string) (spotify.ID, error) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		return "", fmt.Errorf("getting userID: %w", err)
	}
	playlist, err := c.spotify.CreatePlaylistForUser(ctx, userID, name, "", true, false)
	if err != nil {
		return "", fmt.Errorf("creating playlist: %w", err)
	}
	return playlist.ID, nil
}

func (c Client) getNewTracksOfArtist(ctx context.Context, artist string) ([]spotify.SimpleTrack, error) {
	query := fmt.Sprintf("artist:%v tag:new", artist)
	res, err := c.spotify.Search(ctx, query, spotify.SearchTypeAlbum)
	if err != nil {
		return nil, fmt.Errorf("searching new album of %v: %w", artist, err)
	}
	var newTracks []spotify.SimpleTrack
	albums := res.Albums
	if albums == nil {
		return newTracks, nil
	}
	for _, album := range albums.Albums {
		tracks, err := c.getTracksOfAlbum(ctx, album.ID)
		if err != nil {
			return nil, fmt.Errorf("getting tracks of  album %v: %w", album.Name, err)
		}
		newTracks = append(newTracks, tracks...)
	}
	return newTracks, nil
}

func (c Client) getTracksOfAlbum(ctx context.Context, albumID spotify.ID) ([]spotify.SimpleTrack, error) {
	page, err := c.spotify.GetAlbumTracks(ctx, albumID)
	if err != nil {
		return nil, err
	}
	return page.Tracks, nil
}

func (c Client) addTracksToPlaylist(ctx context.Context, playlistID spotify.ID, tracks []spotify.SimpleTrack) error {
	trackIDs := make([]spotify.ID, len(tracks))
	for index, track := range tracks {
		trackIDs[index] = track.ID
	}
	_, err := c.spotify.AddTracksToPlaylist(ctx, playlistID, trackIDs...)
	return err
}

func (c Client) getUserID(ctx context.Context) (string, error) {
	user, err := c.spotify.CurrentUser(ctx)
	if err != nil {
		return "", fmt.Errorf("getting current user: %w", err)
	}
	return user.ID, nil
}

func newSpotifyClient(ctx context.Context, token *oauth2.Token) *spotify.Client {
	httpClient := spotifyauth.New().Client(ctx, token)
	return spotify.New(httpClient)
}
