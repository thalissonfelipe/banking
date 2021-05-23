package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (t Transfer) ListTransfers(ctx context.Context, accountID vos.ID) ([]entities.Transfer, error) {
	transfers, err := t.repository.GetTransfers(ctx, accountID)
	if err != nil {
		return nil, entities.ErrInternalError
	}

	return transfers, nil
}
