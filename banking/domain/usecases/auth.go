package usecases

import (
	"context"
	"errors"
)

var (
	// ErrInvalidCredentials occurs when the user provide invalid credentials.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUnauthorized occurs when the user is not authenticated.
	ErrUnauthorized = errors.New("invalid token")
)

type Auth interface {
	Autheticate(ctx context.Context, cpf, secret string) (token string, err error)
}
