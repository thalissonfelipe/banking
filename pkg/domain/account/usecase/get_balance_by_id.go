package usecase

import (
	"context"
)

func (a Account) GetAccountBalanceByID(ctx context.Context, accountID string) (int, error) {
	balance, err := a.repository.GetBalanceByID(ctx, accountID)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
