package account

import (
	"context"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

type Usecase interface {
	ListAccounts(ctx context.Context) ([]entities.Account, error)
	GetAccountBalanceByID(ctx context.Context, accountID vos.AccountID) (int, error)
	CreateAccount(ctx context.Context, input CreateAccountInput) (*entities.Account, error)
	GetAccountByID(ctx context.Context, accountID vos.AccountID) (entities.Account, error)
	GetAccountByCPF(ctx context.Context, cpf vos.CPF) (entities.Account, error)
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
