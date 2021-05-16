package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type StubRepository struct {
	accounts []entities.Account
	err      error
}

func (s StubRepository) GetAccounts(ctx context.Context) ([]entities.Account, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.accounts, nil
}

func (s StubRepository) GetBalanceByID(ctx context.Context, id string) (int, error) {
	return 0, nil
}

func TestListAccounts(t *testing.T) {
	ctx := context.Background()

	t.Run("should return a list of accounts", func(t *testing.T) {
		acc := entities.NewAccount("Piter", "12345678", "123.456.789-00")
		repo := StubRepository{accounts: []entities.Account{acc}, err: nil}
		usecase := Account{&repo}
		expected := []entities.Account{acc}
		result, err := usecase.ListAccounts(ctx)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("should return an error if something went wrong on repository", func(t *testing.T) {
		repo := StubRepository{accounts: nil, err: errors.New("failed to fetch accounts")}
		usecase := Account{&repo}
		result, err := usecase.ListAccounts(ctx)

		assert.Nil(t, result)
		assert.NotNil(t, err)
	})
}
