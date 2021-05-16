package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type StubRepository struct {
	transfers []entities.Transfer
}

func (s StubRepository) GetTransfers(ctx context.Context) ([]entities.Transfer, error) {
	return s.transfers, nil
}
