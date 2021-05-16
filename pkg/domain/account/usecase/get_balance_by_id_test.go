package usecase

import (
	"context"
	"errors"
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
		expected := 0
		result, err := usecase.GetAccountBalanceByID(ctx, acc.ID)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("should return an error if something went wrong on repository", func(t *testing.T) {
		repo := StubRepository{accounts: nil, err: errors.New("could not fetch account balance")}
		usecase := Account{&repo}
		result, err := usecase.GetAccountBalanceByID(ctx, entities.NewAccountID())

		assert.Zero(t, result)
		assert.NotNil(t, err)
	})
}
