package secrets

import (
	"context"
	"testing"
	"time"

	"github.com/hown3d/kainotomia/test/k8s"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	corev1 "k8s.io/api/core/v1"
)

func Test_createOrUpdateSecret(t *testing.T) {
	type args struct {
		userID string
		token  oauth2.Token
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		storedSecrets []*corev1.Secret
	}{
		{
			name: "create a new secret",
			args: args{
				userID: "12345",
				token: oauth2.Token{
					AccessToken:  "123",
					TokenType:    "123",
					RefreshToken: "123",
					Expiry:       time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr:       false,
			storedSecrets: []*corev1.Secret{},
		},
		// {
		// 	name: "update secret",
		// 	args: args{
		// 		userID: "12345",
		// 		token: oauth2.Token{
		// 			AccessToken:  "123",
		// 			TokenType:    "123",
		// 			RefreshToken: "123",
		// 			Expiry:       time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		// 		},
		// 	},
		// 	wantErr: false,
		// 	storedSecrets: []*corev1.Secret{
		// 		{
		// 			ObjectMeta: v1.ObjectMeta{
		// 				Namespace: "default",
		// 				// name == userID
		// 				Name: "12345",
		// 			},
		// 		},
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientset := k8s.NewFakeClienSet(k8s.ToObjects(tt.storedSecrets)...)
			secretsClient := clientset.CoreV1().Secrets("default")
			secret, err := createOrUpdateSecret(context.Background(), tt.args.userID, tt.args.token, secretsClient)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			token, err := tokenFromSecret(secret)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.args.token, token)
		})
	}
}
