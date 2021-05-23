package transfer

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

type UseCase interface {
	ListTransfers(ctx context.Context, accountID vos.ID) ([]entities.Transfer, error)
	CreateTransfer(ctx context.Context, input CreateTransferInput) error
}

type CreateTransferInput struct {
	AccountOriginID      vos.ID
	AccountDestinationID vos.ID
	Amount               int
}

func NewTransferInput(accOriginID, accDestID vos.ID, amount int) CreateTransferInput {
	return CreateTransferInput{
		AccountOriginID:      accOriginID,
		AccountDestinationID: accDestID,
		Amount:               amount,
	}
}
