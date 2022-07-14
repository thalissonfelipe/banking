package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
)

func TestHash_Hash(t *testing.T) {
	h := Hash{}

	hash, err := h.Hash("12345678")
	assert.NoError(t, err)
	assert.NotNil(t, hash)
}

func TestHash_CompareHashAndSecret(t *testing.T) {
	h := New()

	testCases := []struct {
		name    string
		secret  []byte
		wantErr error
	}{
		{
			name:    "should return no error when the secret is valid",
			secret:  []byte("12345678"),
			wantErr: nil,
		},
		{
			name:    "should return an error when the secret is not valid",
			secret:  []byte("87654321"),
			wantErr: usecases.ErrInvalidCredentials,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := h.Hash("12345678")
			require.NoError(t, err)

			err = h.CompareHashAndSecret(hash, tt.secret)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
