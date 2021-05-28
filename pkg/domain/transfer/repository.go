package transfer

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

type Repository interface {
	GetTransfers(ctx context.Context, id vos.ID) ([]entities.Transfer, error)
	CreateTransfer(ctx context.Context, transfer *entities.Transfer) error
}
