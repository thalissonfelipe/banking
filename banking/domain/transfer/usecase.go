package transfer

import (
	"context"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

type Usecase interface {
	ListTransfers(context.Context, vos.AccountID) ([]entities.Transfer, error)
	PerformTransfer(context.Context, PerformTransferInput) error
}

type PerformTransferInput struct {
	AccountOriginID      vos.AccountID
	AccountDestinationID vos.AccountID
	Amount               int
}

func NewPerformTransferInput(accOriginID, accDestID vos.AccountID, amount int) PerformTransferInput {
	return PerformTransferInput{
		AccountOriginID:      accOriginID,
		AccountDestinationID: accDestID,
		Amount:               amount,
	}
}
