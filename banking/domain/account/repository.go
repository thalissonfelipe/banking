package account

import (
	"context"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

type Repository interface {
	ListAccounts(context.Context) ([]entity.Account, error)
	GetAccountBalanceByID(context.Context, vos.AccountID) (int, error)
	CreateAccount(context.Context, *entity.Account) error
	GetAccountByCPF(context.Context, vos.CPF) (entity.Account, error)
	GetAccountByID(context.Context, vos.AccountID) (entity.Account, error)
}
