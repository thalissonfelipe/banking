package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountRepository_GetAccountByID(t *testing.T) {
	t.Run("should get account by id successfully", func(t *testing.T) {
		db := pgDocker.DB
		r := NewRepository(db)
		ctx := context.Background()

		defer dockertest.TruncateTables(ctx, db)

		account := entities.NewAccount("name", testdata.GetValidCPF(), testdata.GetValidSecret())

		err := r.CreateAccount(ctx, &account)
		require.NoError(t, err)

		got, err := r.GetAccountByID(ctx, account.ID)
		require.NoError(t, err)

		assert.Equal(t, account.ID, got.ID)
		assert.Equal(t, account.Name, got.Name)
		assert.Equal(t, account.CPF, got.CPF)
		assert.Equal(t, account.Balance, got.Balance)
		assert.Equal(t, account.CreatedAt, got.CreatedAt)
	})

	t.Run("should return an error if account does not exist", func(t *testing.T) {
		db := pgDocker.DB
		r := NewRepository(db)
		ctx := context.Background()

		account, err := r.GetAccountByID(ctx, vos.NewAccountID())
		assert.ErrorIs(t, err, entities.ErrAccountDoesNotExist)
		assert.Zero(t, account)
	})
}
