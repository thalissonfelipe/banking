package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer"
)

func (t Transfer) CreateTransfer(ctx context.Context, input transfer.CreateTransferInput) error {
	accOrigin, err := t.accountUsecase.GetAccountByID(ctx, input.AccountOriginID)
	if err != nil {
		if errors.Is(err, entities.ErrAccountDoesNotExist) {
			return fmt.Errorf("accout origin does not exist: %w", err)
		}

		return entities.ErrInternalError
	}

	_, err = t.accountUsecase.GetAccountByID(ctx, input.AccountDestinationID)
	if err != nil {
		if errors.Is(err, entities.ErrAccountDoesNotExist) {
			return entities.ErrAccountDestinationDoesNotExist
		}

		return entities.ErrInternalError
	}

	if (accOrigin.Balance - input.Amount) < 0 {
		return entities.ErrInsufficientFunds
	}

	transfer := entities.NewTransfer(
		input.AccountOriginID,
		input.AccountDestinationID,
		input.Amount,
	)

	err = t.repository.CreateTransfer(ctx, &transfer)
	if err != nil {
		return entities.ErrInternalError
	}

	return nil
}
