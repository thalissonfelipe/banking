package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewToken(t *testing.T) {
	token, err := NewToken("account_id")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestIsValidToken(t *testing.T) {
	token, err := NewToken("account_id")
	require.NoError(t, err)

	err = IsValidToken(token)
	assert.NoError(t, err)
}

func TestGetIDFromToken(t *testing.T) {
	want := "account_id"

	token, err := NewToken(want)
	require.NoError(t, err)

	got := GetIDFromToken(token)
	assert.Equal(t, want, got)
}
