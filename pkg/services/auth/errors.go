package auth

import "errors"

var (
	ErrSecretDoesNotMatch = errors.New("secret does not match")
	ErrUnauthorized       = errors.New("invalid token")
)
