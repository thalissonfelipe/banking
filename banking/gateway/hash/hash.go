package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hash struct{}

func (h Hash) Hash(secret string) ([]byte, error) {
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("generating hash from password: %w", err)
	}

	return hashedSecret, nil
}

func (h Hash) CompareHashAndSecret(hashedSecret, secret []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedSecret, secret)
	if err != nil {
		return fmt.Errorf("comparing hash and password: %w", err)
	}

	return nil
}
