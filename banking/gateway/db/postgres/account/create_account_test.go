package account

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountRepository_CreateAccount(t *testing.T) {
	t.Parallel()

	t.Run("should create an account successfully", func(t *testing.T) {
		t.Parallel()

		db := dockertest.NewDB(t, t.Name())
		r := NewRepository(db)
		ctx := context.Background()

		want, err := entity.NewAccount("name", testdata.CPF().String(), testdata.Secret().String())
		require.NoError(t, err)

		err = r.CreateAccount(ctx, &want)
		require.NoError(t, err)
		assert.True(t, want.CreatedAt.Before(time.Now()))

		acc, err := r.GetAccountByID(ctx, want.ID)
		require.NoError(t, err)

		assert.Equal(t, want.ID, acc.ID)
		assert.Equal(t, want.Name, acc.Name)
		assert.Equal(t, want.CPF, acc.CPF)
		assert.Equal(t, want.Balance, acc.Balance)
		assert.Equal(t, want.Secret, acc.Secret)
	})

	t.Run("should return an error if account already exists", func(t *testing.T) {
		t.Parallel()

		db := dockertest.NewDB(t, t.Name())
		r := NewRepository(db)
		ctx := context.Background()

		acc, err := entity.NewAccount("name", testdata.CPF().String(), testdata.Secret().String())
		require.NoError(t, err)

		err = r.CreateAccount(ctx, &acc)
		require.NoError(t, err)

		err = r.CreateAccount(ctx, &acc)
		assert.ErrorIs(t, err, entity.ErrAccountAlreadyExists)
	})
}
