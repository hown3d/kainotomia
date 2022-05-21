package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/hown3d/kainotomia/pkg/k8s/secrets"
	"github.com/hown3d/kainotomia/pkg/random"
	"github.com/hown3d/kainotomia/pkg/spotify"
	"golang.org/x/oauth2"

	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var (
	state = random.String(10)
	// codeVerifier must be minimum 32 bytes long
	codeVerifier = generateCodeVerifier(32)
)

type Server struct {
	secretsClient corev1.SecretInterface
	spotifyAuth   spotify.Authenticator
	addr          string
}

func NewServer(port int, secretsClient corev1.SecretInterface, spotifyAuth spotify.Authenticator) *Server {
	return &Server{
		addr:          fmt.Sprintf(":%d", port),
		secretsClient: secretsClient,
		spotifyAuth:   spotifyAuth,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", s.authHandler)
	mux.HandleFunc("/", s.redirectHandler)
	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) redirectHandler(w http.ResponseWriter, r *http.Request) {
	url, err := authURL(s.spotifyAuth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Add("Cache-Control", "no-cache")
	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}

func (s *Server) authHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("recived auth")
	st := r.FormValue("state")
	if st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	token, err := s.spotifyAuth.Token(r.Context(), st, r, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	spotifyClient, err := spotify.NewClient(r.Context(), token, s.spotifyAuth)
	err = secrets.StoreToken(r.Context(), spotifyClient.UserID, *token, s.secretsClient)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error storing token, %w", err)
		return
	}

	w.Header().Add("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
}

func authURL(authenticator spotify.Authenticator) (string, error) {
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
