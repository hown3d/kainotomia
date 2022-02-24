package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hown3d/kainotomia/pkg/file"
	"github.com/hown3d/kainotomia/pkg/random"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const (
	redirectURL = "http://localhost:8080/callback"
	clientID    = "2893fb6c7e2c44b8bd85a6a0a8d14033"
)

var state = random.String(20)

func StartHTTPCallbackServer() error {
	authenticator := newSpotifyAuthenticator()
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", authHandler(authenticator))
	mux.HandleFunc("/", redirectHandler(authenticator))
	return http.ListenAndServe(":8080", mux)
}

func redirectHandler(authenticator *spotifyauth.Authenticator) func(w http.ResponseWriter, r *http.Request) {
	url := authURL(authenticator)
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("recived redirect")
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}

func authHandler(authenticator *spotifyauth.Authenticator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("recived auth")
		st := r.FormValue("state")
		if st != state {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, state)
		}
		token, err := authenticator.Token(r.Context(), st, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		file, err := file.ForToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
		}

		err = persistToken(file, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func newSpotifyAuthenticator() *spotifyauth.Authenticator {
	return spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURL),
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithScopes(spotifyauth.ScopePlaylistModifyPublic),
	)
}

func persistToken(w io.Writer, token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("marshaling token", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("writing token: %w", err)
	}
	return nil
}

func authURL(authenticator *spotifyauth.Authenticator) string {
	return authenticator.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", random.String(20)),
	)
}
