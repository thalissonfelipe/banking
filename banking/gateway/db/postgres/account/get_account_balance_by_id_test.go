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

func TestAccountRepository_GetAccountBalanceByID(t *testing.T) {
	t.Parallel()

	t.Run("should get balance successfully", func(t *testing.T) {
		t.Parallel()

		db := dockertest.NewDB(t, t.Name())
		r := NewRepository(db)
		ctx := context.Background()

		wantBalance := 100

		want, err := entity.NewAccount("name", testdata.CPF().String(), testdata.Secret().String())
		require.NoError(t, err)
		want.Balance = wantBalance

		err = r.CreateAccount(ctx, &want)
		require.NoError(t, err)

		balance, err := r.GetAccountBalanceByID(ctx, want.ID)
		require.NoError(t, err)
		assert.Equal(t, wantBalance, balance)
	})

	t.Run("should return an error if account does not exist", func(t *testing.T) {
		t.Parallel()

		db := dockertest.NewDB(t, t.Name())
		r := NewRepository(db)
		ctx := context.Background()

		balance, err := r.GetAccountBalanceByID(ctx, vos.NewAccountID())
		assert.ErrorIs(t, err, entity.ErrAccountNotFound)
		assert.Zero(t, balance)
	})
}
