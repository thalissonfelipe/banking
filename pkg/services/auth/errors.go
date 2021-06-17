package auth

import "errors"

var (
	// ErrInvalidCredentials occurs when the user provide invalid credentials.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUnauthorized occurs when the user is not authenticated.
	ErrUnauthorized = errors.New("invalid token")
)
