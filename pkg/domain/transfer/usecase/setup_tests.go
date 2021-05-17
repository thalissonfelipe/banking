package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type StubRepository struct {
	transfers []entities.Transfer
	err       error
}

func (s StubRepository) GetTransfers(ctx context.Context) ([]entities.Transfer, error) {
	if s.err != nil {
		return nil, entities.ErrInternalError
	}
	return s.transfers, nil
}
