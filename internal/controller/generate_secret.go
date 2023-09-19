package controller

import (
	"crypto/rand"
	"math/big"

	corev1 "k8s.io/api/core/v1"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generates a random string of the given length
func generateRandomString(length int32) (string, error) {
	charlen := big.NewInt(int64(len(charset)))

	// select random characters from the charset
	secret := make([]byte, length)
	for i := range secret {
		charIndex, err := rand.Int(rand.Reader, charlen)
		if err != nil {
			return "", err
		}
		secret[i] = charset[charIndex.Int64()]
	}
	return string(secret), nil
}

// sets random data for the given secret
func setRandomData(secret *corev1.Secret, dataLength int32) error {
	randomString, err := generateRandomString(dataLength)
	if err != nil {
		return err
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	secret.Data["value"] = []byte(randomString)
	return nil
}
