package hash

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
)

var _ encrypter.Encrypter = (*Hash)(nil)

type Hash struct{}

func New() *Hash {
	return &Hash{}
}

func (Hash) Hash(secret string) ([]byte, error) {
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("generating hash from password: %w", err)
	}

	return hashedSecret, nil
}

func (Hash) CompareHashAndSecret(hashedSecret, secret []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedSecret, secret)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return usecases.ErrInvalidCredentials
		}

		return fmt.Errorf("comparing hash and password: %w", err)
	}

	return nil
}
