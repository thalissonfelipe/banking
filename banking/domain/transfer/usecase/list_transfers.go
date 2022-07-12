package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func (t Transfer) ListTransfers(ctx context.Context, accountID vos.AccountID) ([]entities.Transfer, error) {
	transfers, err := t.repository.GetTransfers(ctx, accountID)
	if err != nil {
		return nil, entities.ErrInternalError
	}

	return transfers, nil
}
