package auth

import (
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"context"
)

type key int

const SpotifyTokenKey key = 1337

func SpotifyAuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	newCtx := context.WithValue(ctx, SpotifyTokenKey, token)
	return newCtx, nil
}
