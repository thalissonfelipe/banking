package transfer

import (
	"context"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

type Repository interface {
	GetTransfers(ctx context.Context, id vos.AccountID) ([]entities.Transfer, error)
	CreateTransfer(ctx context.Context, transfer *entities.Transfer) error
}
