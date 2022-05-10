package service

import (
	"context"

	"github.com/hown3d/kainotomia/pkg/config"
	"github.com/hown3d/kainotomia/pkg/job"
	"github.com/hown3d/kainotomia/pkg/spotify"
	kainotomia "github.com/hown3d/kainotomia/proto/kainotomia/v1alpha1"
	"github.com/hown3d/kainotomia/service/auth"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	kubeclient    kubernetes.Interface
	cronjobClient job.Client
	namespace     string
}

func New(cfg config.Config, kubeclient kubernetes.Interface) (Service, error) {
	cjClient := job.NewClient(kubeclient, cfg.Namespace, cfg.JobImage)
	return Service{
		kubeclient:    kubeclient,
		cronjobClient: cjClient,
	}, nil
}

func (s Service) CreatePlaylist(ctx context.Context, req *kainotomia.CreatePlaylistRequest) (*kainotomia.CreatePlaylistResponse, error) {
	oauthToken := tokenFromContext(ctx)
	client, err := spotify.NewClient(oauthToken)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	playlistID, err := client.CreatePlaylist(ctx, req.GetName(), req.GetArtists())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	err = s.cronjobClient.Create(ctx, client.UserID, playlistID.String(), req.GetArtists(), oauthToken)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	return &kainotomia.CreatePlaylistResponse{
		Id: string(playlistID),
	}, nil
}

func (s Service) DeletePlaylist(ctx context.Context, req *kainotomia.DeletePlaylistRequest) (*kainotomia.DeletePlaylistResponse, error) {
	panic("todo: implement")
}

func (s Service) TriggerUpdate(ctx context.Context, req *kainotomia.TriggerUpdateRequest) (*kainotomia.TriggerUpdateResponse, error) {
	panic("todo: implement")
}

func spotifyClientFromToken(token *oauth2.Token) (spotify.Client, error) {
	return spotify.NewClient(token)
}

func tokenFromContext(ctx context.Context) *oauth2.Token {
	return ctx.Value(auth.SpotifyTokenKey).(*oauth2.Token)
}
