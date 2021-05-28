package auth

import "github.com/thalissonfelipe/banking/pkg/domain/account"

type Auth struct {
	accountUsecase account.Usecase
	encrypter      account.Encrypter
}

func NewAuth(accUsecase account.Usecase, encrypter account.Encrypter) *Auth {
	return &Auth{
		accountUsecase: accUsecase,
		encrypter:      encrypter,
	}
}
