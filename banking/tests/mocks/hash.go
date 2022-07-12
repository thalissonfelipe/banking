package mocks

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
	"github.com/thalissonfelipe/banking/banking/domain/entities"
)

var _ encrypter.Encrypter = (*HashMock)(nil)

type HashMock struct {
	Err error
}

func (s HashMock) Hash(secret string) ([]byte, error) {
	if s.Err != nil {
		return nil, entities.ErrInternalError
	}

	return []byte(generateRandomSecret(len(secret))), nil
}

func (s HashMock) CompareHashAndSecret(hashedSecret, secret []byte) error {
	return s.Err
}

func generateRandomSecret(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}
