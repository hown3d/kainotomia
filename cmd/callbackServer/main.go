package main

import (
	"log"

	"github.com/hown3d/kainotomia/pkg/spotify/auth"
)

func main() {
	log.Println("Serving on :8080")
	err := auth.StartHTTPCallbackServer()
	if err != nil {
		log.Fatal(err)
	}
}
