package account

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type UseCase interface {
	ListAccounts(ctx context.Context) ([]entities.Account, error)
	GetAccountBalanceByID(ctx context.Context, accountID string) (int, error)
	CreateAccount(ctx context.Context, input CreateAccountInput) (*entities.Account, error)
	GetAccountByID(ctx context.Context, accountID string) (*entities.Account, error)
}

type CreateAccountInput struct {
	Name   string
	CPF    string
	Secret string
}

func NewCreateAccountInput(name, cpf, secret string) CreateAccountInput {
	return CreateAccountInput{
		Name:   name,
		CPF:    cpf,
		Secret: secret,
	}
}
