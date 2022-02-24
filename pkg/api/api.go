package api

import (
	"context"

	"github.com/hown3d/kainotomia/pkg/api/auth"
	"github.com/hown3d/kainotomia/pkg/spotify"
	kainotomia "github.com/hown3d/kainotomia/proto/kainotomia/v1alpha1"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type API struct {
}

func New() API {
	return API{}
}

func spotifyClientFromToken(ctx context.Context) spotify.Client {
	token := ctx.Value(auth.SpotifyTokenKey)
	return spotify.New(&oauth2.Token{
		AccessToken: token.(string),
	})
}

func (a API) CreatePlaylist(ctx context.Context, req *kainotomia.CreatePlaylistRequest) (*kainotomia.CreatePlaylistResponse, error) {
	client := spotifyClientFromToken(ctx)
	playlistID, err := client.CreatePlaylist(ctx, req.GetName(), req.GetArtists())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	return &kainotomia.CreatePlaylistResponse{
		Id: string(playlistID),
	}, nil
}
