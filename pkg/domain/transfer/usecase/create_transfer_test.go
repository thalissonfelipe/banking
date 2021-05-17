package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer"
)

func TestCreateTransfer(t *testing.T) {
	ctx := context.Background()

	t.Run("should perform a transfer successfully", func(t *testing.T) {
		acc1 := entities.NewAccount("Pedro", "123.456.789-00", "12345678")
		acc1.Balance = 100
		acc2 := entities.NewAccount("Maria", "123.456.789-01", "12345678")

		transferInput := transfer.CreateTransferInput{
			AccountOriginID:      acc1.ID,
			AccountDestinationID: acc2.ID,
			Amount:               100,
		}

		repo := StubRepository{}
		accUseCase := StubAccountUseCase{
			accounts: []entities.Account{acc1, acc2},
		}
		usecase := NewTransfer(&repo, accUseCase)

		err := usecase.CreateTransfer(ctx, transferInput)

		assert.Nil(t, err)
		assert.Len(t, repo.transfers, 1)
	})
}
