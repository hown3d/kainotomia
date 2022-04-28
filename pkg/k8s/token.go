package k8s

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreClientV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// StoreTokenSecret creates a new secret inside kubernetes which stores the oauth2 token.
func StoreTokenInSecret(ctx context.Context, secretsClient coreClientV1.SecretInterface, userID string, token *oauth2.Token) error {
	secretData, err := tokenToMap(token)
	if err != nil {
		return fmt.Errorf("parsing oauth token to map: %w", err)
	}
	_, err = secretsClient.Create(ctx, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: userID,
		},
		Data: secretData,
	}, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating secret for %v: %w", userID, err)
	}
	return nil
}

const (
	AccessTokenSecretKey  string = "SPOTIFY_ACCESS_TOKEN"
	RefreshTokenSecretKey string = "SPOTIFY_REFRESH_TOKEN"
	ExpirySecretKey       string = "SPOTIFY_EXPIRY_TIME"
)

type tokenMap map[string][]byte

func tokenToMap(token *oauth2.Token) (tokenMap, error) {
	m := make(tokenMap)
	m[AccessTokenSecretKey] = []byte(token.AccessToken)
	m[RefreshTokenSecretKey] = []byte(token.RefreshToken)
	expiry, err := token.Expiry.MarshalText()
	if err != nil {
		return nil, fmt.Errorf("marshaling expiry time of token: %w", err)
	}
	m[ExpirySecretKey] = expiry
	return m, nil
}
