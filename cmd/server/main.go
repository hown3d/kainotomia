package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/hown3d/kainotomia/pkg/api"
	"github.com/hown3d/kainotomia/pkg/api/auth"

	"github.com/authzed/grpcutil"

	kainotomiapb "github.com/hown3d/kainotomia/proto/kainotomia/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var port *int = flag.Int("port", 8080, "port to listen on")

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", *port))
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %w", port, err)
	}

	srv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(auth.SpotifyAuthFunc)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(auth.SpotifyAuthFunc)),
	)
	reflection.Register(grpcutil.NewAuthlessReflectionInterceptor(srv))
	api := api.New()
	kainotomiapb.RegisterKainotomiaServiceServer(srv, api)

	log.Printf("serving on %v", lis.Addr().String())
	if err = srv.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
