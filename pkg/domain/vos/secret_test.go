package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSecret(t *testing.T) {
	testCases := []struct {
		name           string
		secret         string
		expectedSecret string
		expectedErr    error
	}{
		{"secret must be valid", "aZ1234Ds", "aZ1234Ds", nil},
		{"secret without uppercase characters", "az1234ds", "", ErrInvalidSecret},
		{"secret without lowercase characters", "AZ1234DS", "", ErrInvalidSecret},
		{"secret without numbers", "azHJKLds", "", ErrInvalidSecret},
		{"secret less than 8 characters", "aZ1234D", "", ErrInvalidSecret},
		{"secret with more than 20 characters", "aZ1234DsaZ1234DsERty0", "", ErrInvalidSecret},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			secret, err := NewSecret(tt.secret)

			assert.Equal(t, tt.expectedSecret, secret.String())
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
