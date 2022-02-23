package auth

import (
	"log"
	"net/http"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const (
	redirectURL = "http://localhost:8080/callback"
	state       = "test123"
)

func StartHTTPCallbackServer(ch chan<- *oauth2.Token) error {
	authenticator := newSpotifyAuthenticator()
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", authHandler(ch, authenticator))
	mux.HandleFunc("/", redirect(authenticator))
	return http.ListenAndServe(":8080", mux)
}

func redirect(authenticator *spotifyauth.Authenticator) func(w http.ResponseWriter, r *http.Request) {
	url := authenticator.AuthURL(state)
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}

func authHandler(ch chan<- *oauth2.Token, authenticator *spotifyauth.Authenticator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		st := r.FormValue("state")
		if st != state {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, state)
		}
		token, err := authenticator.Token(r.Context(), st, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("State mismatch: %s != %s\n", st, state)
		}

		ch <- token
	}
}

func newSpotifyAuthenticator() *spotifyauth.Authenticator {
	return spotifyauth.New(spotifyauth.WithRedirectURL(redirectURL))
}
