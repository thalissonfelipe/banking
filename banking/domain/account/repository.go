package account

import (
	"context"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

type Repository interface {
	ListAccounts(context.Context) ([]entities.Account, error)
	GetAccountBalanceByID(context.Context, vos.AccountID) (int, error)
	CreateAccount(context.Context, *entities.Account) error
	GetAccountByCPF(context.Context, vos.CPF) (entities.Account, error)
	GetAccountByID(context.Context, vos.AccountID) (entities.Account, error)
}
