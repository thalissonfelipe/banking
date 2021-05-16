package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (t Transfer) ListTransfers(ctx context.Context) ([]entities.Transfer, error) {
	transfers, err := t.repository.GetTransfers(ctx)
	return transfers, err
}
