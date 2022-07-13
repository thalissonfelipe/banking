package usecase

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func (a Account) GetAccountByID(ctx context.Context, accountID vos.AccountID) (entity.Account, error) {
	acc, err := a.repository.GetAccountByID(ctx, accountID)
	if err != nil {
		return entity.Account{}, fmt.Errorf("getting account by id: %w", err)
	}

	return acc, nil
}
