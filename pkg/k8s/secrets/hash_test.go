package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_annotateSecretWithHash(t *testing.T) {
	type args struct {
		secret *corev1.Secret
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "check if secret has annotation",
			args: args{
				secret: &corev1.Secret{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := annotateSecretWithHash(tt.args.secret)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Contains(t, tt.args.secret.Annotations, hashAnnotation)
		})
	}
}

func Test_secretChanged(t *testing.T) {
	type args struct {
		newSecret     *corev1.Secret
		currentSecret *corev1.Secret
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "same secrets",
			args: args{
				newSecret: &corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name: "secret",
					},
				},
				currentSecret: &corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name: "secret",
					},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "different secrets",
			args: args{
				newSecret: &corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name: "new secret",
					},
				},
				currentSecret: &corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name: "current secret",
					},
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := annotateSecretWithHash(tt.args.currentSecret)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			got, err := secretChanged(tt.args.newSecret, tt.args.currentSecret)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
