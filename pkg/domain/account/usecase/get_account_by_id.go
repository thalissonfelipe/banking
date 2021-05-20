package usecase

import (
	"context"
	"errors"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (a Account) GetAccountByID(ctx context.Context, accountID string) (*entities.Account, error) {
	acc, err := a.repository.GetAccountByID(ctx, accountID)
	if err != nil {
		if !errors.Is(err, entities.ErrAccountDoesNotExist) {
			return nil, entities.ErrInternalError
		}
		return nil, err
	}

	return acc, nil
}
