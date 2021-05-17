package transfer

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type UseCase interface {
	ListTransfers(ctx context.Context, accountID string) ([]entities.Transfer, error)
}
