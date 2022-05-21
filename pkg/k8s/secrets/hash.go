package secrets

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	corev1 "k8s.io/api/core/v1"
)

const hashAnnotation string = "kainotomia.io/last-applied-hash"

func annotateSecretWithHash(secret *corev1.Secret, hash uint64) error {
	secret.Annotations[hashAnnotation] = string(hash)
	return nil
}

func secretChanged(new *corev1.Secret, current *corev1.Secret) (bool, error) {
	newHash, err := secretHash(new)
	if err != nil {
		return false, err
	}
	currentHash := current.Annotations[hashAnnotation]
	return currentHash == string(newHash), nil
}

func secretHash(secret *corev1.Secret) (uint64, error) {
	hash, err := hashstructure.Hash(secret, hashstructure.FormatV2, nil)
	if err != nil {
		return 0, fmt.Errorf("hashing secret: %w", err)
	}
	return hash, nil
}
