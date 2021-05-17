package usecase

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer"
)

func (t Transfer) CreateTransfer(ctx context.Context, input transfer.CreateTransferInput) error {
	accountOriginBalance, err := t.accountUseCase.GetAccountBalanceByID(ctx, input.AccountOriginID)
	if err != nil {
		return err
	}

	accDestination, err := t.accountUseCase.GetAccountByID(ctx, input.AccountDestinationID)
	if err != nil {
		return err
	}
	if accDestination == nil {
		return entities.ErrAccountDestinationDoesNotExist
	}

	if (accountOriginBalance - input.Amount) < 0 {
		return entities.ErrInsufficientFunds
	}

	transfer := entities.NewTransfer(
		input.AccountOriginID,
		input.AccountDestinationID,
		input.Amount,
	)
	err = t.repository.UpdateBalance(ctx, transfer)

	return err
}
