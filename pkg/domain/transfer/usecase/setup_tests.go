package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type StubRepository struct {
	transfers []entities.Transfer
	err       error
}

func (s StubRepository) GetTransfers(ctx context.Context, id string) ([]entities.Transfer, error) {
	if s.err != nil {
		return nil, entities.ErrInternalError
	}

	var transfers []entities.Transfer
	for _, tr := range s.transfers {
		if tr.AccountOriginID == id {
			transfers = append(transfers, tr)
		}
	}

	return transfers, nil
}

func (s StubRepository) UpdateBalance(ctx context.Context, transfer entities.Transfer) error {
	return nil
}
