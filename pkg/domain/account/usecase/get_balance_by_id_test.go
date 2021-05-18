package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
)

func TestGetBalanceByAccountID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return a balance by account ID", func(t *testing.T) {
		acc := entities.NewAccount("Piter", "123.456.789-00", "12345678")
		repo := mocks.StubAccountRepository{Accounts: []entities.Account{acc}}
		usecase := NewAccountUseCase(&repo, nil)
		expected := entities.DefaultBalance
		result, err := usecase.GetAccountBalanceByID(ctx, acc.ID)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("should return an error if account does not exist", func(t *testing.T) {
		repo := mocks.StubAccountRepository{Accounts: nil}
		usecase := NewAccountUseCase(&repo, nil)
		result, err := usecase.GetAccountBalanceByID(ctx, entities.NewAccountID())

		assert.Zero(t, result)
		assert.Equal(t, err, entities.ErrAccountDoesNotExist)
	})

	t.Run("should return correct balance when balance is not 0", func(t *testing.T) {
		acc := entities.NewAccount("Piter", "123.456.789-00", "12345678")
		acc.Balance = 100
		repo := mocks.StubAccountRepository{Accounts: []entities.Account{acc}}
		usecase := NewAccountUseCase(&repo, nil)
		expected := 100
		result, err := usecase.GetAccountBalanceByID(ctx, acc.ID)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
}
