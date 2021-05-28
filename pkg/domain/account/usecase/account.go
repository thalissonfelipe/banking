package usecase

import "github.com/thalissonfelipe/banking/pkg/domain/account"

type Account struct {
	repository account.Repository
	encrypter  account.Encrypter
}

func NewAccountUsecase(repo account.Repository, enc account.Encrypter) *Account {
	return &Account{
		repository: repo,
		encrypter:  enc,
	}
}
