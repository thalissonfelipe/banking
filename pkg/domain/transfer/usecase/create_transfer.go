package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer"
)

func (t Transfer) CreateTransfer(ctx context.Context, input transfer.CreateTransferInput) error {
	accOrigin, err := t.accountUseCase.GetAccountByID(ctx, input.AccountOriginID)
	if err != nil {
		return entities.ErrInternalError
	}
	if accOrigin == nil {
		return entities.ErrAccountDoesNotExist
	}

	accDestination, err := t.accountUseCase.GetAccountByID(ctx, input.AccountDestinationID)
	if err != nil {
		return entities.ErrInternalError
	}
	if accDestination == nil {
		return entities.ErrAccountDestinationDoesNotExist
	}

	if (accOrigin.Balance - input.Amount) < 0 {
		return entities.ErrInsufficientFunds
	}

	transfer := entities.NewTransfer(
		input.AccountOriginID,
		input.AccountDestinationID,
		input.Amount,
	)
	err = t.repository.UpdateBalance(ctx, transfer)
	if err != nil {
		return entities.ErrInternalError
	}

	return nil
}
