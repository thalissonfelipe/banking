package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func (a Auth) Autheticate(ctx context.Context, input AuthenticateInput) (string, error) {
	cpf, err := vos.NewCPF(input.CPF)
	if err != nil {
		return "", fmt.Errorf("invalid cpf: %w", err)
	}

	acc, err := a.accountUsecase.GetAccountByCPF(ctx, cpf)
	if err != nil {
		if errors.Is(err, entities.ErrInternalError) {
			return "", fmt.Errorf("account does not exist: %w", err)
		}

		return "", ErrInvalidCredentials
	}

	hashedSecret := []byte(acc.Secret.String())
	secret := []byte(input.Secret)

	err = a.encrypter.CompareHashAndSecret(hashedSecret, secret)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := NewToken(acc.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}
