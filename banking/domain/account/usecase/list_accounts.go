package usecase

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
)

func (a Account) ListAccounts(ctx context.Context) ([]entity.Account, error) {
	accounts, err := a.repository.ListAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing accounts: %w", err)
	}

	return accounts, nil
}
