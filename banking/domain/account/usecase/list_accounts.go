package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
)

func (a Account) ListAccounts(ctx context.Context) ([]entities.Account, error) {
	accounts, err := a.repository.GetAccounts(ctx)
	if err != nil {
		return nil, entities.ErrInternalError
	}

	return accounts, nil
}
