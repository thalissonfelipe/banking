package transfer

import (
	"context"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

type Repository interface {
	ListTransfers(context.Context, vos.AccountID) ([]entities.Transfer, error)
	PerformTransfer(context.Context, *entities.Transfer) error
}
