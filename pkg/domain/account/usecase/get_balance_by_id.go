package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (a Account) GetAccountBalanceByID(ctx context.Context, accountID string) (int, error) {
	accExists, err := a.repository.GetAccountByID(ctx, accountID)
	if err != nil {
		return 0, entities.ErrInternalError
	}
	if accExists == nil {
		return 0, entities.ErrAccountDoesNotExist
	}

	balance, err := a.repository.GetBalanceByID(ctx, accountID)
	if err != nil {
		return 0, entities.ErrInternalError
	}

	return balance, nil
}
