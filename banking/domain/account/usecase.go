package account

import (
	"context"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

type Usecase interface {
	ListAccounts(context.Context) ([]entities.Account, error)
	GetAccountBalanceByID(context.Context, vos.AccountID) (int, error)
	CreateAccount(context.Context, *entities.Account) error
	GetAccountByID(context.Context, vos.AccountID) (entities.Account, error)
	GetAccountByCPF(context.Context, vos.CPF) (entities.Account, error)
}

type CreateAccountInput struct {
	Name   string
	CPF    vos.CPF
	Secret vos.Secret
}

func NewCreateAccountInput(name string, cpf vos.CPF, secret vos.Secret) CreateAccountInput {
	return CreateAccountInput{
		Name:   name,
		CPF:    cpf,
		Secret: secret,
	}
}
