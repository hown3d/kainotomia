package service

import (
	"context"

	"github.com/hown3d/kainotomia/pkg/config"
	"github.com/hown3d/kainotomia/pkg/k8s/job"
	"github.com/hown3d/kainotomia/pkg/spotify"
	kainotomia "github.com/hown3d/kainotomia/proto/kainotomia/v1alpha1"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	kubeclient    kubernetes.Interface
	cronjobClient job.Client
	spotifyAuth   spotify.Authenticator
	namespace     string
}

func New(cfg config.Config, kubeclient kubernetes.Interface, spotifyAuth spotify.Authenticator) (Service, error) {
	cjClient := job.NewClient(kubeclient, cfg.Namespace, cfg.JobImage)
	return Service{
		kubeclient:    kubeclient,
		cronjobClient: cjClient,
		spotifyAuth:   spotifyAuth,
	}, nil
}

func (s Service) CreatePlaylist(ctx context.Context, req *kainotomia.CreatePlaylistRequest) (*kainotomia.CreatePlaylistResponse, error) {
	token := tokenFromContext(ctx)
	client, err := spotify.NewClient(ctx, token, s.spotifyAuth)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	playlistID, err := client.CreatePlaylist(ctx, req.GetName(), req.GetArtists())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	err = s.cronjobClient.Create(ctx, client.UserID, playlistID.String(), req.GetArtists(), token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
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

func tokenFromContext(ctx context.Context) *oauth2.Token {
	return ctx.Value(spotifyTokenKey).(*oauth2.Token)
}
