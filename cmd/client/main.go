package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hown3d/kainotomia/pkg/spotify/auth"
	"golang.org/x/oauth2"
)

var url *string = flag.String("url", "http://localhost:8081", "")

func main() {
	flag.Parse()
	ch := make(chan *oauth2.Token)
	log.Println("Serving on :8081")
	go func() {
		err := auth.StartHTTPCallbackServer(*url, ch)
		if err != nil {
			log.Fatal(err)
		}
	}()
	token := <-ch
	fmt.Print(token.AccessToken)
}
