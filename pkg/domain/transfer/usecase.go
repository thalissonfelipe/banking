package transfer

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type UseCase interface {
	ListTransfers(ctx context.Context) ([]entities.Transfer, error)
	// TODO: CreateTransfer
}
