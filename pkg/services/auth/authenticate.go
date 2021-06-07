package auth

import (
	"context"
	"errors"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (a Auth) Autheticate(ctx context.Context, input AuthenticateInput) (string, error) {
	cpf, err := vos.NewCPF(input.CPF)
	if err != nil {
		return "", err
	}

	acc, err := a.accountUsecase.GetAccountByCPF(ctx, cpf)
	if err != nil {
		if errors.Is(err, entities.ErrInternalError) {
			return "", err
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

	return token, err
}
