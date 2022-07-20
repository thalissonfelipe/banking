package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountRepository_ListAccounts(t *testing.T) {
	t.Parallel()

	db := dockertest.NewDB(t, t.Name())
	r := NewRepository(db)
	ctx := context.Background()

	accounts, err := r.ListAccounts(ctx)
	require.NoError(t, err)
	assert.Len(t, accounts, 0)

	want, err := entity.NewAccount("name", testdata.CPF().String(), testdata.Secret().String())
	require.NoError(t, err)

	err = r.CreateAccount(ctx, &want)
	require.NoError(t, err)

	accounts, err = r.ListAccounts(ctx)
	require.NoError(t, err)

	assert.Len(t, accounts, 1)
	assert.Equal(t, want.ID, accounts[0].ID)
	assert.Equal(t, want.Name, accounts[0].Name)
	assert.Equal(t, want.CPF.String(), accounts[0].CPF.String())
	assert.Equal(t, want.Balance, accounts[0].Balance)
}
