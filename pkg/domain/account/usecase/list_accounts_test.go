package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
)

func TestListAccounts(t *testing.T) {
	ctx := context.Background()

	t.Run("should return a list of accounts", func(t *testing.T) {
		acc := entities.NewAccount("Piter", "123.456.789-00", "12345678")
		repo := mocks.StubAccountRepository{Accounts: []entities.Account{acc}, Err: nil}
		usecase := NewAccountUseCase(&repo, nil)
		expected := []entities.Account{acc}
		result, err := usecase.ListAccounts(ctx)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("should return an error if something went wrong on repository", func(t *testing.T) {
		repo := mocks.StubAccountRepository{Accounts: nil, Err: errors.New("failed to fetch accounts")}
		usecase := NewAccountUseCase(&repo, nil)
		result, err := usecase.ListAccounts(ctx)

		assert.Nil(t, result)
		assert.NotNil(t, err)
	})
}
