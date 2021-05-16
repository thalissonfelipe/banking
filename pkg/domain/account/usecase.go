package account

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type UseCase interface {
	ListAccounts(ctx context.Context) ([]entities.Account, error)
	GetAccountBalanceByID(ctx context.Context, accountID string) (int, error)
	// TODO: CreateAccount
}
