package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (a Account) GetAccountByID(ctx context.Context, accountID vos.AccountID) (*entities.Account, error) {
	acc, err := a.repository.GetAccountByID(ctx, accountID)
	if err == nil {
		return acc, nil
	}

	if errors.Is(err, entities.ErrAccountDoesNotExist) {
		return nil, fmt.Errorf("account does not exist: %w", err)
	}

	return nil, entities.ErrInternalError
}
