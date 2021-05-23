package auth

import (
	"context"
)

func (a Auth) Autheticate(ctx context.Context, input AuthenticateInput) (string, error) {
	acc, err := a.accountUsecase.GetAccountByCPF(ctx, input.CPF)
	if err != nil {
		return "", err
	}

	hashedSecret := []byte(acc.Secret.String())
	secret := []byte(input.Secret)

	err = a.encrypter.CompareHashAndSecret(hashedSecret, secret)
	if err != nil {
		return "", ErrSecretDoesNotMatch
	}

	token, err := NewToken(acc.ID.String())

	return token, err
}
