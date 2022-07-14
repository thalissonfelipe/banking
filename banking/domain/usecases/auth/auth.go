package auth

import (
	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
)

//go:generate moq -pkg auth -out repository_mock.gen.go ../../entity AccountRepository:RepositoryMock
//go:generate moq -pkg auth -out encrypter_mock.gen.go ../../encrypter Encrypter

var _ usecases.Auth = (*Auth)(nil)

type Auth struct {
	usecase   usecases.Account
	encrypter encrypter.Encrypter
}

func NewAuth(usecase usecases.Account, encrypter encrypter.Encrypter) *Auth {
	return &Auth{
		usecase:   usecase,
		encrypter: encrypter,
	}
}
