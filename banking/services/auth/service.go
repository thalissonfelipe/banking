package auth

import (
	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/services"
)

//go:generate moq -pkg auth -out repository_mock.gen.go ../../domain/entity AccountRepository:RepositoryMock
//go:generate moq -pkg auth -out encrypter_mock.gen.go ../../domain/encrypter Encrypter

var _ services.Auth = (*Auth)(nil)

type Auth struct {
	accountUsecase usecases.Account
	encrypter      encrypter.Encrypter
}

func NewAuth(accUsecase usecases.Account, encrypter encrypter.Encrypter) *Auth {
	return &Auth{
		accountUsecase: accUsecase,
		encrypter:      encrypter,
	}
}
