package usecase

import (
	"context"
	"errors"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (a Account) GetAccountByID(ctx context.Context, accountID vos.ID) (*entities.Account, error) {
	acc, err := a.repository.GetAccountByID(ctx, accountID)
	if err != nil {
		if !errors.Is(err, entities.ErrAccountDoesNotExist) {
			return nil, entities.ErrInternalError
		}
		return nil, err
	}

	return acc, nil
}
