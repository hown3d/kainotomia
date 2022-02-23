package main

import (
	"log"

	"github.com/hown3d/kainotomia/pkg/spotify/auth"
	"golang.org/x/oauth2"
)

func main() {
	tokenChan := make(chan *oauth2.Token)
	err := auth.StartHTTPCallbackServer(tokenChan)
	if err != nil {
		log.Fatal(err)
	}
}
