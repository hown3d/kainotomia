package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/hown3d/kainotomia/pkg/config"
	"github.com/hown3d/kainotomia/pkg/k8s"
	spotify_auth "github.com/hown3d/kainotomia/pkg/spotify/auth"
	"github.com/hown3d/kainotomia/service"
	"github.com/hown3d/kainotomia/service/auth"

	"github.com/authzed/grpcutil"

	kainotomiapb "github.com/hown3d/kainotomia/proto/kainotomia/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port     *int = flag.Int("port", 8080, "port to listen on")
	authPort *int = flag.Int("auth-port", 8081, "auth port for token retrieval port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", *port))
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", port, err)
	}

	srv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(auth.SpotifyAuthFunc)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(auth.SpotifyAuthFunc)),
	)
	reflection.Register(grpcutil.NewAuthlessReflectionInterceptor(srv))

	cfg := config.Parse()
	kubeclient, err := k8s.NewClientSet()
	if err != nil {
		log.Fatal(fmt.Errorf("creating new kubernetes client set: %w", err))
	}
	service, err := service.New(cfg, kubeclient)
	if err != nil {
		log.Fatal(err)
	}
	kainotomiapb.RegisterKainotomiaServiceServer(srv, service)

	go func() {
		srvUrl := fmt.Sprintf("http://0.0.0.0:%v", *authPort)
		log.Printf("serving token server on %v", srvUrl)
		err := spotify_auth.StartHTTPCallbackServer(srvUrl, kubeclient.CoreV1().Secrets(cfg.Namespace))
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Printf("serving on %v", lis.Addr().String())
	if err = srv.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
