package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountRepository_GetAccountByID(t *testing.T) {
	t.Parallel()

	t.Run("should get account by id successfully", func(t *testing.T) {
		t.Parallel()

		db := dockertest.NewDB(t, t.Name())
		r := NewRepository(db)
		ctx := context.Background()

		want, err := entity.NewAccount("name", testdata.CPF().String(), testdata.Secret().String())
		require.NoError(t, err)

		err = r.CreateAccount(ctx, &want)
		require.NoError(t, err)

		got, err := r.GetAccountByID(ctx, want.ID)
		require.NoError(t, err)

		assert.Equal(t, want.ID, got.ID)
		assert.Equal(t, want.Name, got.Name)
		assert.Equal(t, want.CPF, got.CPF)
		assert.Equal(t, want.Balance, got.Balance)
		assert.Equal(t, want.CreatedAt, got.CreatedAt)
	})

	t.Run("should return an error if account does not exist", func(t *testing.T) {
		t.Parallel()

		db := dockertest.NewDB(t, t.Name())
		r := NewRepository(db)
		ctx := context.Background()

		account, err := r.GetAccountByID(ctx, vos.NewAccountID())
		assert.ErrorIs(t, err, entity.ErrAccountNotFound)
		assert.Zero(t, account)
	})
}
