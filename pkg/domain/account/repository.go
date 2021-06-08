package account

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

type Repository interface {
	GetAccounts(ctx context.Context) ([]entities.Account, error)
	GetBalanceByID(ctx context.Context, id vos.AccountID) (int, error)
	CreateAccount(ctx context.Context, account *entities.Account) error
	GetAccountByCPF(ctx context.Context, cpf vos.CPF) (*entities.Account, error)
	GetAccountByID(ctx context.Context, id vos.AccountID) (*entities.Account, error)
}
