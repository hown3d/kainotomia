package secrets

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreClientV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	accessTokenSecretKey  string = "TOKEN_ACCESS_TOKEN"
	refreshTokenSecretKey string = "TOKEN_REFRESH_TOKEN"
	expirySecretKey       string = "TOKEN_EXPIRY_TIME"
	typeSecretKey         string = "TOKEN_TYPE"
)

// StoreToken creates a new secret inside kubernetes which stores the oauth2 token.
func StoreToken(ctx context.Context, userID string, token oauth2.Token, secretsClient coreClientV1.SecretInterface) (*corev1.Secret, error) {
	return createOrUpdateSecret(ctx, userID, token, secretsClient)
}

func LoadToken(ctx context.Context, userID string, secretsClient coreClientV1.SecretInterface) (oauth2.Token, error) {
	getOpts := metav1.GetOptions{}
	secret, err := secretsClient.Get(ctx, userID, getOpts)
	if err != nil {
		return oauth2.Token{}, fmt.Errorf("getting token secret %w", err)
	}
	return tokenFromSecret(secret)
}

func createOrUpdateSecret(ctx context.Context, userID string, token oauth2.Token, client coreClientV1.SecretInterface) (*corev1.Secret, error) {
	// check if secret exists
	getOpts := metav1.GetOptions{}
	secret, err := client.Get(ctx, userID, getOpts)
	if err != nil {
		if errors.IsNotFound(err) {
			secret, err := secretFromToken(userID, token)
			if err != nil {
				return nil, fmt.Errorf("generating secret from token: %w", err)
			}
			return createSecret(ctx, secret, client)

		} else {
			return nil, fmt.Errorf("getting token secret: %w", err)
		}
	}

	// secret exists, so check if the hashes are still the same
	newSecret, err := secretFromToken(userID, token)
	if err != nil {
		return nil, fmt.Errorf("generating secret from token: %w", err)
	}
	changed, err := secretChanged(newSecret, secret)
	if err != nil {
		return nil, fmt.Errorf("checking if secret changed: %w", err)
	}
	if changed {
		return updateSecret(ctx, newSecret, client)
	}
	return secret, nil
}

func secretFromToken(userID string, token oauth2.Token) (*corev1.Secret, error) {
	secretData, err := tokenToMap(token)
	if err != nil {
		return nil, fmt.Errorf("parsing token into map: %w", err)
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: userID,
		},
		Data: secretData,
	}
	err = annotateSecretWithHash(secret)
	if err != nil {
		return nil, fmt.Errorf("annotating token secret with hash: %w", err)
	}
	return secret, nil
}

func tokenFromSecret(secret *corev1.Secret) (oauth2.Token, error) {
	tokenData := secret.Data
	expiryTime := new(time.Time)
	err := expiryTime.UnmarshalText(tokenData[expirySecretKey])
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

func createSecret(ctx context.Context, secret *corev1.Secret, client coreClientV1.SecretInterface) (*corev1.Secret, error) {
	new, err := client.Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("creating new secret for token: %w", err)
	}
	return new, nil
}

func updateSecret(ctx context.Context, secret *corev1.Secret, client coreClientV1.SecretInterface) (*corev1.Secret, error) {
	updated, err := client.Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("creating new secret for token: %w", err)
	}
	return updated, nil
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
