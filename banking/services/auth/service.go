package auth

import (
	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
	"github.com/thalissonfelipe/banking/banking/services"
)

//go:generate moq -pkg auth -out repository_mock.gen.go ../../domain/account Repository
//go:generate moq -pkg auth -out encrypter_mock.gen.go ../../domain/encrypter Encrypter

var _ services.Auth = (*Auth)(nil)

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
