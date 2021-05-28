package usecase

import (
	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/encrypter"
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
