package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHash_Hash(t *testing.T) {
	h := Hash{}

	hash, err := h.Hash("12345678")
	assert.NoError(t, err)
	assert.NotNil(t, hash)
}

func TestHash_CompareHashAndSecret(t *testing.T) {
	h := Hash{}

	testCases := []struct {
		name    string
		secret  []byte
		hash    func(t *testing.T) []byte
		wantErr error
	}{
		{
			name:   "should return no error when the secret is valid",
			secret: []byte("12345678"),
			hash: func(t *testing.T) []byte {
				hash, err := h.Hash("12345678")
				require.NoError(t, err)

				return hash
			},
			wantErr: nil,
		},
		{
			name:   "should return an error when the secret is not valid",
			secret: []byte("87654321"),
			hash: func(t *testing.T) []byte {
				hash, err := h.Hash("12345678")
				require.NoError(t, err)

				return hash
			},
			wantErr: bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := h.CompareHashAndSecret(tt.hash(t), tt.secret)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
