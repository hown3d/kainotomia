package secrets

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	corev1 "k8s.io/api/core/v1"
)

const hashAnnotation string = "kainotomia.io/last-applied-hash"

func annotateSecretWithHash(secret *corev1.Secret) error {
	hash, err := secretHash(secret)
	if err != nil {
		return err
	}
	if secret.Annotations == nil {
		secret.Annotations = make(map[string]string)
	}
	secret.Annotations[hashAnnotation] = fmt.Sprint(hash)
	return nil
}

func secretChanged(newSecret *corev1.Secret, currentSecret *corev1.Secret) (bool, error) {
	newHash, err := secretHash(newSecret)
	if err != nil {
		return false, err
	}
	currentHash := currentSecret.Annotations[hashAnnotation]
	return currentHash == fmt.Sprint(newHash), nil
}

func secretHash(secret *corev1.Secret) (uint64, error) {
	hash, err := hashstructure.Hash(*secret, hashstructure.FormatV2, nil)
	if err != nil {
		return 0, fmt.Errorf("hashing secret: %w", err)
	}
	return hash, nil
}
