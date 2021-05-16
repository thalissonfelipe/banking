package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

type StubRepository struct {
	accounts []entities.Account
}

func (s StubRepository) GetAccounts(ctx context.Context) ([]entities.Account, error) {
	return s.accounts, nil
}

func (s StubRepository) GetBalanceByID(ctx context.Context, id string) (int, error) {
	return 0, nil
}

func TestListAccounts(t *testing.T) {
	acc := entities.NewAccount("Piter", "12345678", "123.456.789-00")
	repo := StubRepository{accounts: []entities.Account{acc}}
	usecase := Account{&repo}
	ctx := context.Background()

	t.Run("should return a list of accounts", func(t *testing.T) {
		expected := []entities.Account{acc}
		result, err := usecase.ListAccounts(ctx)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	})
}
