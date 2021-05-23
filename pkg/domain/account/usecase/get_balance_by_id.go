package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (a Account) GetAccountBalanceByID(ctx context.Context, accountID vos.ID) (int, error) {
	balance, err := a.repository.GetBalanceByID(ctx, accountID)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
