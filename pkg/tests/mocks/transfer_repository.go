package mocks

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

var _ transfer.Repository = (*TransferRepositoryMock)(nil)

type TransferRepositoryMock struct {
	Transfers []entities.Transfer
	Err       error
}

func (s TransferRepositoryMock) GetTransfers(ctx context.Context, id vos.ID) ([]entities.Transfer, error) {
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

func (s *TransferRepositoryMock) CreateTransfer(ctx context.Context, transfer *entities.Transfer) error {
	if s.Err != nil {
		return entities.ErrInternalError
	}

	s.Transfers = append(s.Transfers, *transfer)

	return nil
}
