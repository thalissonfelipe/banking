package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewToken(t *testing.T) {
	t.Parallel()

	token, err := NewToken("account_id")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestIsTokenValid(t *testing.T) {
	t.Parallel()

	token, err := NewToken("account_id")
	require.NoError(t, err)

	err = IsTokenValid(token)
	assert.NoError(t, err)

	err = IsTokenValid("random_token")
	assert.Error(t, err)
}

func TestGetIDFromToken(t *testing.T) {
	t.Parallel()

	const want = "account_id"

	token, err := NewToken(want)
	require.NoError(t, err)

	got := GetAccountIDFromToken(token)
	assert.Equal(t, want, got)
}
