package auth

import (
	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/encrypter"
)

type Auth struct {
	accountUsecase account.UseCase
	encrypter      encrypter.Encrypter
}

func NewAuth(accUsecase account.UseCase, encrypter encrypter.Encrypter) *Auth {
	return &Auth{
		accountUsecase: accUsecase,
		encrypter:      encrypter,
	}
}
