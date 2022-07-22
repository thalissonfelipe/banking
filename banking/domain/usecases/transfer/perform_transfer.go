package transfer

import (
	"context"
	"errors"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
)

func (u Usecase) PerformTransfer(ctx context.Context, input usecases.PerformTransferInput) error {
	accOrigin, err := u.accountUsecase.GetAccountByID(ctx, input.AccountOriginID)
	if err != nil {
		return fmt.Errorf("getting origin account by id: %w", err)
	}

	_, err = u.accountUsecase.GetAccountByID(ctx, input.AccountDestinationID)
	if err != nil {
		if errors.Is(err, entity.ErrAccountNotFound) {
			return entity.ErrAccountDestinationNotFound
		}

		return fmt.Errorf("getting destination account by id: %w", err)
	}

	transfer, err := entity.NewTransfer(
		input.AccountOriginID,
		input.AccountDestinationID,
		input.Amount,
		accOrigin.Balance,
	)
	if err != nil {
		return fmt.Errorf("creating transfer: %w", err)
	}

	err = u.repository.PerformTransfer(ctx, &transfer)
	if err != nil {
		return fmt.Errorf("creating transfer: %w", err)
	}

	return nil
}
