package auth

import (
	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
)

//go:generate moq -pkg auth -out repository_mock.gen.go ../../entity AccountRepository:RepositoryMock
//go:generate moq -pkg auth -out encrypter_mock.gen.go ../../encrypter Encrypter
//go:generate moq -pkg auth -out service_mock.gen.go . Service

var _ usecases.Auth = (*Auth)(nil)

type Auth struct {
	usecase   usecases.Account
	encrypter encrypter.Encrypter
	service   Service
}

func NewAuth(usecase usecases.Account, encrypter encrypter.Encrypter, service Service) *Auth {
	return &Auth{
		usecase:   usecase,
		encrypter: encrypter,
		service:   service,
	}
}

type Service interface {
	NewToken(accountID string) (token string, err error)
}
