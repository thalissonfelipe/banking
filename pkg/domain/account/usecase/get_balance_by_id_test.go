package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func TestGetBalanceByAccountID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return a balance by account ID", func(t *testing.T) {
		acc := entities.NewAccount("Piter", "12345678", "123.456.789-00")
		repo := StubRepository{accounts: []entities.Account{acc}, err: nil}
		usecase := Account{&repo}
		expected := entities.DefaultBalance
		result, err := usecase.GetAccountBalanceByID(ctx, acc.ID)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("should return an error if account does not exist", func(t *testing.T) {
		repo := StubRepository{accounts: nil, err: entities.ErrAccountDoesNotExist}
		usecase := Account{&repo}
		result, err := usecase.GetAccountBalanceByID(ctx, entities.NewAccountID())

		assert.Zero(t, result)
		assert.Equal(t, err, entities.ErrAccountDoesNotExist)
	})

	t.Run("should return correct balance when balance is not 0", func(t *testing.T) {
		acc := entities.NewAccount("Piter", "12345678", "123.456.789-00")
		acc.Balance = 100
		repo := StubRepository{accounts: []entities.Account{acc}, err: nil}
		usecase := Account{&repo}
		expected := 100
		result, err := usecase.GetAccountBalanceByID(ctx, acc.ID)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
}
