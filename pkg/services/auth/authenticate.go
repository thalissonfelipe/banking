package auth

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (a Auth) Autheticate(ctx context.Context, input AuthenticateInput) (string, error) {
	acc, err := a.accountUsecase.GetAccountByCPF(ctx, input.CPF)
	if err != nil {
		return "", err
	}
	if acc == nil {
		return "", entities.ErrAccountDoesNotExist
	}

	hashedSecret := []byte(acc.Secret.String())
	secret := []byte(input.Secret)

	err = a.encrypter.CompareHashAndSecret(hashedSecret, secret)
	if err != nil {
		return "", ErrSecretDoesNotMatch
	}

	token, err := newToken(acc.ID)

	return token, err
}
