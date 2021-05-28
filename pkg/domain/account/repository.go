package account

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

type Repository interface {
	GetAccounts(ctx context.Context) ([]entities.Account, error)
	GetBalanceByID(ctx context.Context, id vos.ID) (int, error)
	PostAccount(ctx context.Context, account *entities.Account) error
	GetAccountByCPF(ctx context.Context, cpf vos.CPF) (*entities.Account, error)
	GetAccountByID(ctx context.Context, id vos.ID) (*entities.Account, error)
}
