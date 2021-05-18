package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
)

func TestCreateTransfer(t *testing.T) {
	ctx := context.Background()
	acc1 := entities.NewAccount("Pedro", "123.456.789-00", "12345678")
	acc2 := entities.NewAccount("Maria", "123.456.789-01", "12345678")

	t.Run("should perform a transfer successfully", func(t *testing.T) {
		acc1.Balance = 100
		transferInput := transfer.CreateTransferInput{
			AccountOriginID:      acc1.ID,
			AccountDestinationID: acc2.ID,
			Amount:               100,
		}

		repo := mocks.StubTransferRepository{}
		accUseCase := mocks.StubAccountUseCase{
			Accounts: []entities.Account{acc1, acc2},
		}
		usecase := NewTransfer(&repo, accUseCase)

		err := usecase.CreateTransfer(ctx, transferInput)

		assert.Nil(t, err)
		assert.Len(t, repo.Transfers, 1)
	})

	t.Run("should return an error if accOrigin does not have sufficient funds", func(t *testing.T) {
		acc1.Balance = 0
		transferInput := transfer.CreateTransferInput{
			AccountOriginID:      acc1.ID,
			AccountDestinationID: acc2.ID,
			Amount:               100,
		}

		repo := mocks.StubTransferRepository{}
		accUseCase := mocks.StubAccountUseCase{
			Accounts: []entities.Account{acc1, acc2},
		}
		usecase := NewTransfer(&repo, accUseCase)

		err := usecase.CreateTransfer(ctx, transferInput)

		assert.Equal(t, err, entities.ErrInsufficientFunds)
		assert.Len(t, repo.Transfers, 0)
	})

	t.Run("should return an error if accountUseCase fails", func(t *testing.T) {
		transferInput := transfer.CreateTransferInput{
			AccountOriginID:      acc1.ID,
			AccountDestinationID: acc2.ID,
			Amount:               100,
		}

		repo := mocks.StubTransferRepository{}
		accUseCase := mocks.StubAccountUseCase{Err: errors.New("failed to get account origin balance")}
		usecase := NewTransfer(&repo, accUseCase)

		err := usecase.CreateTransfer(ctx, transferInput)

		assert.NotNil(t, err)
		assert.Len(t, repo.Transfers, 0)
	})

	t.Run("should return an error if account destination does not exist", func(t *testing.T) {
		acc1.Balance = 100
		transferInput := transfer.CreateTransferInput{
			AccountOriginID:      acc1.ID,
			AccountDestinationID: "undefined",
			Amount:               100,
		}

		repo := mocks.StubTransferRepository{Err: errors.New("account dest does not exist")}
		accUseCase := mocks.StubAccountUseCase{
			Accounts: []entities.Account{acc1, acc2},
		}
		usecase := NewTransfer(&repo, accUseCase)

		err := usecase.CreateTransfer(ctx, transferInput)

		assert.Equal(t, err, entities.ErrAccountDestinationDoesNotExist)
		assert.Len(t, repo.Transfers, 0)
	})
}
