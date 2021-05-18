package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretIsValid(t *testing.T) {
	testCases := []struct {
		name     string
		secret   string
		expected bool
	}{
		{"secret must be valid", "aZ1234Ds", true},
		{"secret without uppercase characters", "az1234ds", false},
		{"secret without lowercase characters", "AZ1234DS", false},
		{"secret without numbers", "azHJKLds", false},
		{"secret less than 8 characters", "aZ1234D", false},
		{"secret with more than 20 characters", "aZ1234DsaZ1234DsERty0", false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			secret := Secret{value: tt.secret}
			result := secret.IsValid()

			assert.Equal(t, tt.expected, result)
		})
	}
}
