package transfer

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type Repository interface {
	GetTransfers(ctx context.Context, id string) ([]entities.Transfer, error)
	UpdateBalance(ctx context.Context, transfer entities.Transfer) error
}
