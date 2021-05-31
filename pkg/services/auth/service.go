package auth

import (
	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/encrypter"
)

type Auth struct {
	accountUsecase account.Usecase
	encrypter      encrypter.Encrypter
}

func NewAuth(accUsecase account.Usecase, encrypter encrypter.Encrypter) *Auth {
	return &Auth{
		accountUsecase: accUsecase,
		encrypter:      encrypter,
	}
}
