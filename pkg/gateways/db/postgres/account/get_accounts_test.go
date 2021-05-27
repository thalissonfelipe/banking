package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
)

func TestRepository_GetAccounts(t *testing.T) {
	db := pgDocker.DB
	r := NewRepository(db)
	ctx := context.Background()

	defer dockertest.TruncateTables(ctx, db)

	accounts, err := r.GetAccounts(ctx)
	assert.NoError(t, err)
	assert.Len(t, accounts, 0)

	acc := entities.NewAccount(
		"Maria",
		vos.NewCPF("123.456.789-00"),
		vos.NewSecret("12345678"),
	)

	err = r.PostAccount(ctx, &acc)
	assert.NoError(t, err)

	accounts, err = r.GetAccounts(ctx)

	assert.NoError(t, err)
	assert.Len(t, accounts, 1)
	assert.Equal(t, acc.ID, accounts[0].ID)
	assert.Equal(t, acc.Name, accounts[0].Name)
	assert.Equal(t, acc.CPF.String(), accounts[0].CPF.String())
	assert.Equal(t, acc.Balance, accounts[0].Balance)
}
