package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/hown3d/kainotomia/pkg/random"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const (
	redirectURL = "http://localhost:8080/callback"
	clientID    = "2893fb6c7e2c44b8bd85a6a0a8d14033"
)

var (
	state = random.String(10)
	// codeVerifier must be minimum 32 bytes long
	codeVerifier = generateCodeVerifier(32)
)

func StartHTTPCallbackServer(ch chan<- *oauth2.Token) error {
	authenticator := newSpotifyAuthenticator()
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", authHandler(authenticator, ch))
	mux.HandleFunc("/", redirectHandler(authenticator))
	return http.ListenAndServe(":8080", mux)
}

func redirectHandler(authenticator *spotifyauth.Authenticator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := authURL(authenticator)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Println("recived redirect")
		w.Header().Add("Cache-Control", "no-cache")
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}

func authHandler(authenticator *spotifyauth.Authenticator, ch chan<- *oauth2.Token) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("recived auth")
		st := r.FormValue("state")
		if st != state {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, state)
		}
		token, err := authenticator.Token(r.Context(), st, r, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// file, err := file.ForToken()
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	log.Fatal(err)
		// }

		// err = persistToken(file, token)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	log.Fatal(err)
		// }
		ch <- token
		w.Header().Add("Cache-Control", "no-cache")
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

func authURL(authenticator *spotifyauth.Authenticator) (string, error) {
	codeChallenge, err := generateCodeChallenge(codeVerifier)
	if err != nil {
		return "", fmt.Errorf("generating code challenge: %w", err)
	}
	return authenticator.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	), nil
}

func generateCodeChallenge(verifier string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(verifier))
	if err != nil {
		return "", fmt.Errorf("hashing verifier string: %w", err)
	}
	encoded := encode(hash.Sum(nil))
	return encoded, nil
}

func encode(msg []byte) string {
	encoded := base64.StdEncoding.EncodeToString(msg)
	encoded = strings.Replace(encoded, "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.Replace(encoded, "=", "", -1)
	return encoded
}

func generateCodeVerifier(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = byte(r.Intn(255))
	}
	return encode(b)
}
