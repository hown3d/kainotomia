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
	UserID  string
}

func NewClient(token *oauth2.Token) (Client, error) {
	c := Client{
		spotify: newSpotifyClient(context.Background(), token),
		logger:  zap.L().Sugar(),
	}
	userID, err := c.getUserID(context.Background())
	if err != nil {
		return c, fmt.Errorf("getting id of user: %w", err)
	}
	c.UserID = userID
	return c, nil
}

func newSpotifyClient(ctx context.Context, token *oauth2.Token) *spotify.Client {
	httpClient := spotifyauth.New().Client(ctx, token)
	return spotify.New(httpClient)
}

func (c Client) getUserID(ctx context.Context) (string, error) {
	user, err := c.spotify.CurrentUser(ctx)
	if err != nil {
		return "", fmt.Errorf("getting current user: %w", err)
	}
	return user.ID, nil
}
