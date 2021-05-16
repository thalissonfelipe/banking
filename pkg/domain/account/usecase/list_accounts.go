package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (a Account) ListAccounts(ctx context.Context) ([]entities.Account, error) {
	accounts, err := a.repository.GetAccounts(ctx)
	return accounts, err
}
