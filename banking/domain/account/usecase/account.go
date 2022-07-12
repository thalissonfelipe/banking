package usecase

import (
	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
)

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
