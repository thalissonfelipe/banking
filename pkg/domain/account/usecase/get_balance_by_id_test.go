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
		expected := 0
		result, err := usecase.GetAccountBalanceByID(ctx, acc.ID)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
}
