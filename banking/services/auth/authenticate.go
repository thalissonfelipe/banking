package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"golang.org/x/crypto/bcrypt"
)

func (a Auth) Autheticate(ctx context.Context, input AuthenticateInput) (string, error) {
	cpf, err := vos.NewCPF(input.CPF)
	if err != nil {
		return "", fmt.Errorf("new cpf: %w", err)
	}

	acc, err := a.accountUsecase.GetAccountByCPF(ctx, cpf)
	if err != nil {
		if errors.Is(err, entities.ErrAccountDoesNotExist) {
			return "", ErrInvalidCredentials
		}

		return "", fmt.Errorf("getting account by cpf: %w", err)
	}

	hashedSecret := []byte(acc.Secret.String())
	secret := []byte(input.Secret)

	err = a.encrypter.CompareHashAndSecret(hashedSecret, secret)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrInvalidCredentials
		}

		return "", fmt.Errorf("hashing secret: %w", err)
	}

	token, err := NewToken(acc.ID.String())
	if err != nil {
		return "", fmt.Errorf("creating token: %w", err)
	}

	return token, nil
}
