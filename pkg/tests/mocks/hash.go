package mocks

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type StubHash struct {
	Err error
}

func (s StubHash) Hash(secret string) ([]byte, error) {
	if s.Err != nil {
		return nil, entities.ErrInternalError
	}

	return []byte(generateRandomSecret(len(secret))), nil
}

func generateRandomSecret(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
