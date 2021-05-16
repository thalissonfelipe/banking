package usecase

import "github.com/thalissonfelipe/banking/pkg/domain/account"

type Account struct {
	repository account.Repository
}

func NewAccountUseCase(repo account.Repository) *Account {
	return &Account{repository: repo}
}
