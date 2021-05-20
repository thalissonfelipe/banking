package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIDFromToken(t *testing.T) {
	t.Run("should return correct id from token", func(t *testing.T) {
		expectedID := "account_id"
		token, _ := newToken(expectedID)
		id, err := GetIDFromToken(token)

		assert.Nil(t, err)
		assert.Equal(t, expectedID, id)
	})
}
