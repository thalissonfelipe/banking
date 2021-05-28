package mocks

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

type StubTransferRepository struct {
	Transfers []entities.Transfer
	Err       error
}

func (s StubTransferRepository) GetTransfers(ctx context.Context, id vos.ID) ([]entities.Transfer, error) {
	if s.Err != nil {
		return nil, entities.ErrInternalError
	}

	var transfers []entities.Transfer
	for _, tr := range s.Transfers {
		if tr.AccountOriginID == id {
			transfers = append(transfers, tr)
		}
	}

	return transfers, nil
}

func (s *StubTransferRepository) CreateTransfer(ctx context.Context, transfer *entities.Transfer) error {
	if s.Err != nil {
		return entities.ErrInternalError
	}
	s.Transfers = append(s.Transfers, *transfer)
	return nil
}
