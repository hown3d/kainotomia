package main

import (
	"flag"
	"log"

	"github.com/hown3d/kainotomia/pkg/spotify/auth"
	"github.com/pkg/browser"
)

var url *string = flag.String("url", "http://localhost:8080", "")

func main() {
	flag.Parse()
	log.Println("Serving on :8080")
	go func() {
		err := auth.StartHTTPCallbackServer()
		if err != nil {
			log.Fatal(err)
		}
	}()
	browser.OpenURL(*url)
}
