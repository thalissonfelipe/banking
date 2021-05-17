package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (a Account) GetAccountByID(ctx context.Context, accountID string) (*entities.Account, error) {
	acc, err := a.repository.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, entities.ErrInternalError
	}

	return acc, nil
}
