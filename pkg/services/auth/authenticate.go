package auth

import (
	"context"

	log "github.com/sirupsen/logrus"
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
		log.WithError(err).Error("unable to compare hashed password with input password")
		return "", ErrSecretDoesNotMatch
	}

	token, err := NewToken(acc.ID.String())
	if err != nil {
		log.WithError(err).Error("unable to create a new token")
		return "", err
	}

	return token, nil
}
