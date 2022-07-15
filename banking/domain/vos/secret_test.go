package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSecret(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		secret  string
		want    string
		wantErr error
	}{
		{"secret must be valid", "aZ1234Ds", "aZ1234Ds", nil},
		{"secret without uppercase characters", "az1234ds", "", ErrInvalidSecret},
		{"secret without lowercase characters", "AZ1234DS", "", ErrInvalidSecret},
		{"secret without numbers", "azHJKLds", "", ErrInvalidSecret},
		{"secret less than 8 characters", "aZ1234D", "", ErrInvalidSecret},
		{"secret with more than 20 characters", "aZ1234DsaZ1234DsERty0", "", ErrInvalidSecret},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			secret, err := NewSecret(tt.secret)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, secret.String())
		})
	}
}
