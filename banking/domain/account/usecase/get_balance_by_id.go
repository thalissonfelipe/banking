package usecase

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func (a Account) GetAccountBalanceByID(ctx context.Context, accountID vos.AccountID) (int, error) {
	balance, err := a.repository.GetAccountBalanceByID(ctx, accountID)
	if err != nil {
		return 0, fmt.Errorf("getting balance by id: %w", err)
	}

	return balance, nil
}
