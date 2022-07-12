package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		name        string
		secret      []byte
		hash        func(t *testing.T) []byte
		expectedErr error
	}{
		{
			name:   "should return nil when the secret is correct",
			secret: []byte("12345678"),
			hash: func(t *testing.T) []byte {
				hash, err := h.Hash("12345678")
				require.NoError(t, err)

				return hash
			},
			expectedErr: nil,
		},
		{
			name:   "should return error when the secret is not correct",
			secret: []byte("87654321"),
			hash: func(t *testing.T) []byte {
				hash, err := h.Hash("12345678")
				require.NoError(t, err)

				return hash
			},
			expectedErr: bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := h.CompareHashAndSecret(tt.hash(t), tt.secret)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
