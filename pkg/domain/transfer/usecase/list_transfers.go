package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (t Transfer) ListTransfers(ctx context.Context, accountID string) ([]entities.Transfer, error) {
	transfers, err := t.repository.GetTransfers(ctx, accountID)
	if err != nil {
		return nil, entities.ErrInternalError
	}

	return transfers, nil
}
