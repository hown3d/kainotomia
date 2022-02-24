package spotify

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type client struct {
	spotify *spotify.Client
}

func New(token *oauth2.Token) client {
	return client{
		spotify: newSpotifyClient(context.Background(), token),
	}
}

func (c client) createPlaylist(ctx context.Context, playlistName string) (spotify.ID, error) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		return "", fmt.Errorf("getting userID: %w", err)
	}
	playlist, err := c.spotify.CreatePlaylistForUser(ctx, userID, playlistName, "", true, false)
	if err != nil {
		return "", fmt.Errorf("creating playlist: %w", err)
	}
	return playlist.ID, nil
}

func (c client) getNewTracksOfArtist(ctx context.Context, artist string) ([]spotify.FullTrack, error) {
	query := fmt.Sprintf("%v tag:new", artist)
	res, err := c.spotify.Search(ctx, query, spotify.SearchTypeTrack)
	if err != nil {
		return nil, fmt.Errorf("searching new album of %v: %w", artist, err)
	}
	trackPage := res.Tracks
	if trackPage == nil {
		return []spotify.FullTrack{}, nil
	}
	return trackPage.Tracks, nil
}

func (c client) getUserID(ctx context.Context) (string, error) {
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
