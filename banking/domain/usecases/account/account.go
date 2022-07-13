package account

import (
	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
)

//go:generate moq -pkg account -out repository_mock.gen.go ../../entity AccountRepository:RepositoryMock
//go:generate moq -pkg account -out encrypter_mock.gen.go ../../encrypter Encrypter

var _ usecases.Account = (*Account)(nil)

type Account struct {
	repository entity.AccountRepository
	encrypter  encrypter.Encrypter
}

func NewAccountUsecase(repo entity.AccountRepository, enc encrypter.Encrypter) *Account {
	return &Account{
		repository: repo,
		encrypter:  enc,
	}
}
