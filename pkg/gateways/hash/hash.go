package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hash struct{}

func (h Hash) Hash(secret string) ([]byte, error) {
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash secret: %w", err)
	}

	return hashedSecret, nil
}

func (h Hash) CompareHashAndSecret(hashedSecret, secret []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedSecret, secret)
	if err != nil {
		return fmt.Errorf("could not compare hash and secret: %w", err)
	}

	return nil
}
