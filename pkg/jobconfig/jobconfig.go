package config

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	coreClientV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// StoreToken creates a new secret inside kubernetes which stores the oauth2 token.
// OwnerReference will be set to the corresponding cronjob, to prevent zombieTokens
func StoreToken(ctx context.Context, secretsClient coreClientV1.SecretInterface, name string, token oauth2.Token, references []v1.OwnerReference) error {
	secretData, err := tokenToMap(token)
	if err != nil {
		return fmt.Errorf("parsing oauth token to map: %w", err)
	}
	_, err = secretsClient.Create(ctx, &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:            name,
			OwnerReferences: references,
		},
		Data: secretData,
	}, v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating secret for %v: %w", name, err)
	}
	return nil
}

type JobConfig struct {
	PlaylistID string
}

func StoreJobConfig(ctx context.Context, cmClient coreClientV1.ConfigMapInterface, name string, jobconfig JobConfig, references []v1.OwnerReference) error {
	data, err := yaml.Marshal(jobconfig)
	if err != nil {
		return fmt.Errorf("marshaling job config: %w", err)
	}
	_, err = cmClient.Create(ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		Data: map[string]string{
			ConfigFileName: string(data),
		},
	}, v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating job config map %v: %w", name, err)
	}
	return nil
}

const (
	AccessTokenSecretKey  string = "SPOTIFY_ACCESS_TOKEN"
	RefreshTokenSecretKey string = "SPOTIFY_REFRESH_TOKEN"
	ExpirySecretKey       string = "SPOTIFY_EXPIRY_TIME"
	PlaylistIDKey         string = "SPOTIFY_PLAYLIST_ID"
	ConfigFilePath        string = "/kainotomia"
	ConfigFileName        string = "config.yml"
)

type tokenMap map[string][]byte

func tokenToMap(token oauth2.Token) (tokenMap, error) {
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
