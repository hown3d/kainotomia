package spotify

import (
	"context"
	"net/http"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

// Authenticator is used to authenticate to spotify via oauth2 flow
type Authenticator interface {
	SpotifyClient(ctx context.Context, token *oauth2.Token, opts ...spotify.ClientOption) SpotifyClient
	AuthURL(state string, opts ...oauth2.AuthCodeOption) string
	Token(ctx context.Context, state string, r *http.Request, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

// Interface Compliance
var _ Authenticator = (*DefaultAuthenticator)(nil)
var _ Authenticator = (*NoOpAuthenticator)(nil)

type DefaultAuthenticator struct {
	authenticator *spotifyauth.Authenticator
}

// Exchange implements Authenticator
func (a *DefaultAuthenticator) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return a.authenticator.Exchange(ctx, code, opts...)
}

func NewDefaultAuthenticator(opts ...spotifyauth.AuthenticatorOption) *DefaultAuthenticator {
	return &DefaultAuthenticator{
		authenticator: spotifyauth.New(opts...),
	}
}

// Token implements Authenticator
func (a *DefaultAuthenticator) Token(ctx context.Context, state string, r *http.Request, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return a.authenticator.Token(ctx, state, r, opts...)
}

func (a *DefaultAuthenticator) SpotifyClient(ctx context.Context, token *oauth2.Token, opts ...spotify.ClientOption) SpotifyClient {
	httpClient := a.authenticator.Client(ctx, token)
	return spotify.New(httpClient, opts...)
}

func (a *DefaultAuthenticator) AuthURL(state string, opts ...oauth2.AuthCodeOption) string {
	return a.authenticator.AuthURL(state, opts...)
}

type NoOpAuthenticator struct{}

// AuthURL implements Authenticator
func (*NoOpAuthenticator) AuthURL(state string, opts ...oauth2.AuthCodeOption) string {
	panic("unimplemented")
}

// Exchange implements Authenticator
func (*NoOpAuthenticator) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	panic("unimplemented")
}

// SpotifyClient implements Authenticator
func (*NoOpAuthenticator) SpotifyClient(ctx context.Context, token *oauth2.Token, opts ...spotify.ClientOption) SpotifyClient {
	panic("unimplemented")
}

// Token implements Authenticator
func (*NoOpAuthenticator) Token(ctx context.Context, state string, r *http.Request, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	panic("unimplemented")
}
