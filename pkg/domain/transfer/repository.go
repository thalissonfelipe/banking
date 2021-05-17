package transfer

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type Repository interface {
	GetTransfers(ctx context.Context, id string) ([]entities.Transfer, error)
}
