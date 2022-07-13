package account

import (
	"context"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

type Usecase interface {
	ListAccounts(context.Context) ([]entity.Account, error)
	GetAccountBalanceByID(context.Context, vos.AccountID) (int, error)
	CreateAccount(context.Context, *entity.Account) error
	GetAccountByID(context.Context, vos.AccountID) (entity.Account, error)
	GetAccountByCPF(context.Context, vos.CPF) (entity.Account, error)
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
