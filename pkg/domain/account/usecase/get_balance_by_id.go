package usecase

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (a Account) GetAccountBalanceByID(ctx context.Context, accountID vos.AccountID) (int, error) {
	balance, err := a.repository.GetBalanceByID(ctx, accountID)
	if err != nil {
		return 0, fmt.Errorf("could not get account balance: %w", err)
	}

	return balance, nil
}
