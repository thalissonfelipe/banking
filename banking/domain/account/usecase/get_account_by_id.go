package usecase

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func (a Account) GetAccountByID(ctx context.Context, accountID vos.AccountID) (entities.Account, error) {
	acc, err := a.repository.GetAccountByID(ctx, accountID)
	if err != nil {
		return entities.Account{}, fmt.Errorf("getting account by id: %w", err)
	}

	return acc, nil
}
