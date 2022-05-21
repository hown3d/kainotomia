package spotify

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// SpotifyClient is the interface to interact with spotify
type SpotifyClient interface {
	CurrentUser(ctx context.Context) (*spotify.PrivateUser, error)
	Search(ctx context.Context, query string, t spotify.SearchType, opts ...spotify.RequestOption) (*spotify.SearchResult, error)
	GetAlbumTracks(ctx context.Context, id spotify.ID, opts ...spotify.RequestOption) (*spotify.SimpleTrackPage, error)
	CreatePlaylistForUser(ctx context.Context, userID string, playlistName string, description string, public bool, collaborative bool) (*spotify.FullPlaylist, error)
	AddTracksToPlaylist(ctx context.Context, playlistID spotify.ID, trackIDs ...spotify.ID) (snapshotID string, err error)
}

type Client struct {
	spotify SpotifyClient
	logger  *zap.SugaredLogger
	UserID  string
}

func NewClient(ctx context.Context, token *oauth2.Token, auth Authenticator) (Client, error) {
	c := Client{
		logger:  zap.L().Sugar(),
		spotify: auth.SpotifyClient(ctx, token),
	}
	userID, err := c.getUserID(context.Background())
	if err != nil {
		return c, fmt.Errorf("getting id of user: %w", err)
	}
	c.UserID = userID
	return c, nil
}

func (c Client) getUserID(ctx context.Context) (string, error) {
	user, err := c.spotify.CurrentUser(ctx)
	if err != nil {
		return "", fmt.Errorf("getting current user: %w", err)
	}
	return user.ID, nil
}
