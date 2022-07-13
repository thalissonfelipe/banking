package usecase

import (
	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
)

//go:generate moq -pkg usecase -out repository_mock.gen.go .. Repository
//go:generate moq -pkg usecase -out encrypter_mock.gen.go ../../encrypter Encrypter

type Account struct {
	repository account.Repository
	encrypter  encrypter.Encrypter
}

func NewAccountUsecase(repo account.Repository, enc encrypter.Encrypter) *Account {
	return &Account{
		repository: repo,
		encrypter:  enc,
	}
}
