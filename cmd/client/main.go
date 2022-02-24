package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hown3d/kainotomia/pkg/spotify/auth"
	"golang.org/x/oauth2"
)

var url *string = flag.String("url", "http://localhost:8080", "")

func main() {
	flag.Parse()
	ch := make(chan *oauth2.Token)
	log.Println("Serving on :8080")
	go func() {
		err := auth.StartHTTPCallbackServer(ch)
		if err != nil {
			log.Fatal(err)
		}
	}()
	token := <-ch
	fmt.Print(token)
}
