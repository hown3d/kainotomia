package service

import (
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"context"
)

type key int

const spotifyTokenKey key = 1337

func SpotifyAuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	newCtx := context.WithValue(ctx, spotifyTokenKey, token)
	return newCtx, nil
}
