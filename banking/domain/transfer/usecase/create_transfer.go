package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/transfer"
)

func (t Transfer) CreateTransfer(ctx context.Context, input transfer.CreateTransferInput) error {
	accOrigin, err := t.accountUsecase.GetAccountByID(ctx, input.AccountOriginID)
	if err != nil {
		return fmt.Errorf("getting origin account by id: %w", err)
	}

	_, err = t.accountUsecase.GetAccountByID(ctx, input.AccountDestinationID)
	if err != nil {
		if errors.Is(err, entities.ErrAccountDoesNotExist) {
			return entities.ErrAccountDestinationDoesNotExist
		}

		return fmt.Errorf("getting destination account by id: %w", err)
	}

	transfer, err := entities.NewTransfer(
		input.AccountOriginID,
		input.AccountDestinationID,
		input.Amount,
		accOrigin.Balance,
	)
	if err != nil {
		return fmt.Errorf("creating transfer: %w", err)
	}

	err = t.repository.CreateTransfer(ctx, &transfer)
	if err != nil {
		return fmt.Errorf("creating transfer: %w", err)
	}

	return nil
}
