package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHash(t *testing.T) {
	t.Run("should create an hash successfully", func(t *testing.T) {
		h := Hash{}
		hash, err := h.Hash("12345678")

		assert.Nil(t, err)
		assert.NotNil(t, hash)
	})
}

func TestCompareHashAndSecret(t *testing.T) {
	h := Hash{}

	testCases := []struct {
		name   string
		secret []byte
		hash   func() []byte
		err    error
	}{
		{
			name:   "should return nil when the secret is correct",
			secret: []byte("12345678"),
			hash: func() []byte {
				hash, _ := h.Hash("12345678")
				return hash
			},
			err: nil,
		},
		{
			name:   "should return error when the secret is not correct",
			secret: []byte("87654321"),
			hash: func() []byte {
				hash, _ := h.Hash("12345678")
				return hash
			},
			err: bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := h.CompareHashAndSecret(tt.hash(), tt.secret)
			assert.Equal(t, tt.err, err)
		})
	}
}
