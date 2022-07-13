package transfer

import (
	"context"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

type Repository interface {
	ListTransfers(context.Context, vos.AccountID) ([]entity.Transfer, error)
	PerformTransfer(context.Context, *entity.Transfer) error
}
