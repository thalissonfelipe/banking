package account

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountRepository_CreateAccount(t *testing.T) {
	t.Run("should create an account successfully", func(t *testing.T) {
		db := pgDocker.DB
		r := NewRepository(db)
		ctx := context.Background()

		defer dockertest.TruncateTables(ctx, db)

		newAccount := entities.NewAccount("name", testdata.GetValidCPF(), testdata.GetValidSecret())

		err := r.CreateAccount(ctx, &newAccount)
		require.NoError(t, err)
		assert.True(t, newAccount.CreatedAt.Before(time.Now()))

		acc, err := r.GetAccountByID(ctx, newAccount.ID)
		require.NoError(t, err)

		assert.Equal(t, newAccount.ID, acc.ID)
		assert.Equal(t, newAccount.Name, acc.Name)
		assert.Equal(t, newAccount.CPF, acc.CPF)
		assert.Equal(t, newAccount.Balance, acc.Balance)
		assert.Equal(t, newAccount.Secret, acc.Secret)
	})

	t.Run("should return an error if account already exists", func(t *testing.T) {
		db := pgDocker.DB
		r := NewRepository(db)
		ctx := context.Background()

		acc := entities.NewAccount("name", testdata.GetValidCPF(), testdata.GetValidSecret())

		err := r.CreateAccount(ctx, &acc)
		require.NoError(t, err)

		err = r.CreateAccount(ctx, &acc)
		assert.ErrorIs(t, err, entities.ErrAccountAlreadyExists)
	})
}
