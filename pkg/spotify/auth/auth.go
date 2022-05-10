package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hown3d/kainotomia/pkg/k8s"
	"github.com/hown3d/kainotomia/pkg/random"
	"github.com/hown3d/kainotomia/pkg/spotify"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"

	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	redirectURL = "http://localhost:8081/callback"
	clientID    = "2893fb6c7e2c44b8bd85a6a0a8d14033"
)

var (
	state = random.String(10)
	// codeVerifier must be minimum 32 bytes long
	codeVerifier = generateCodeVerifier(32)
)

func StartHTTPCallbackServer(rawURL string, secretsClient corev1.SecretInterface) error {
	uri, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	authenticator, err := newSpotifyAuthenticator(uri)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", authHandler(authenticator, secretsClient))
	mux.HandleFunc("/", redirectHandler(authenticator))
	return http.ListenAndServe(fmt.Sprintf(":%v", uri.Port()), mux)
}

func redirectHandler(authenticator *spotifyauth.Authenticator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := authURL(authenticator)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Add("Cache-Control", "no-cache")
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}

func authHandler(authenticator *spotifyauth.Authenticator, secretClient corev1.SecretInterface) func(w http.ResponseWriter, r *http.Request) {
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

		spotifyClient, err := spotify.NewClient(token)
		err = k8s.StoreTokenInSecret(r.Context(), secretClient, spotifyClient.UserID, *token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("error storing token, %w", err)
			return
		}

		w.Header().Add("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
	}
}

func newSpotifyAuthenticator(uri *url.URL) (*spotifyauth.Authenticator, error) {
	redirect, err := uri.Parse("/callback")
	if err != nil {
		return nil, err
	}
	return spotifyauth.New(
		spotifyauth.WithRedirectURL(redirect.String()),
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithScopes(spotifyauth.ScopePlaylistModifyPublic),
	), nil
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
