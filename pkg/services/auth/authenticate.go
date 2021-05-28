package auth

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (a Auth) Autheticate(ctx context.Context, input AuthenticateInput) (string, error) {
	cpf := vos.NewCPF(input.CPF)
	acc, err := a.accountUsecase.GetAccountByCPF(ctx, cpf)
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
