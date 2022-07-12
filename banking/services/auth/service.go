package auth

import (
	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
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

type AuthenticateInput struct {
	CPF    string
	Secret string
}
