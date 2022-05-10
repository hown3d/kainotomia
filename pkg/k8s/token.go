package k8s

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreClientV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	accessTokenSecretKey  string = "TOKEN_ACCESS_TOKEN"
	refreshTokenSecretKey string = "TOKEN_REFRESH_TOKEN"
	expirySecretKey       string = "TOKEN_EXPIRY_TIME"
	typeSecretKey         string = "TOKEN_TYPE"
)

// StoreTokenSecret creates a new secret inside kubernetes which stores the oauth2 token.
func StoreTokenInSecret(ctx context.Context, secretsClient coreClientV1.SecretInterface, userID string, token oauth2.Token) error {
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

func LoadTokenFromSecret(ctx context.Context, secretsClient coreClientV1.SecretInterface, secretName string) (oauth2.Token, error) {
	secret, err := secretsClient.Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return oauth2.Token{}, fmt.Errorf("getting token secret %w", err)
	}
	tokenData := secret.Data
	expiryTime := new(time.Time)
	err = expiryTime.UnmarshalText(tokenData[expirySecretKey])
	if err != nil {
		return oauth2.Token{}, fmt.Errorf("unmarshaling expiry time of token: %w", err)
	}
	return oauth2.Token{
		AccessToken:  string(tokenData[accessTokenSecretKey]),
		RefreshToken: string(tokenData[refreshTokenSecretKey]),
		TokenType:    string(tokenData[typeSecretKey]),
		Expiry:       *expiryTime,
	}, nil
}

type tokenMap map[string][]byte

func tokenToMap(token oauth2.Token) (tokenMap, error) {
	m := make(tokenMap)
	m[accessTokenSecretKey] = []byte(token.AccessToken)
	m[refreshTokenSecretKey] = []byte(token.RefreshToken)
	expiry, err := token.Expiry.MarshalText()
	if err != nil {
		return nil, fmt.Errorf("marshaling expiry time of token: %w", err)
	}
	m[expirySecretKey] = expiry
	m[typeSecretKey] = []byte(token.TokenType)
	return m, nil
}
