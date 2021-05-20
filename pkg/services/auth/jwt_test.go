package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	t.Run("should create a token successfully", func(t *testing.T) {
		accountID := "account_id"
		token, err := newToken(accountID)

		assert.Nil(t, err)
		assert.NotNil(t, token)
	})
}

func TestIsValidToken(t *testing.T) {
	t.Run("should return nil when token is valid", func(t *testing.T) {
		accountID := "account_id"
		token, _ := newToken(accountID)
		err := IsValidToken(token)

		assert.Nil(t, err)
	})
}

func TestGetIDFromToken(t *testing.T) {
	t.Run("should get id from token", func(t *testing.T) {
		accountID := "account_id"
		token, _ := newToken(accountID)
		id := GetIDFromToken(token)

		assert.Equal(t, accountID, id)
	})
}
