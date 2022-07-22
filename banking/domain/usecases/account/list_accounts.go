package account

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
)

func (u Usecase) ListAccounts(ctx context.Context) ([]entity.Account, error) {
	accounts, err := u.repository.ListAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing accounts: %w", err)
	}

	return accounts, nil
}
