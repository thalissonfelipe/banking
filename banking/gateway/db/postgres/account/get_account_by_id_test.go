package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestRepository_GetAccountByID(t *testing.T) {
	db := pgDocker.DB
	r := NewRepository(db)
	ctx := context.Background()

	defer dockertest.TruncateTables(ctx, db)

	randomID := vos.NewAccountID()

	_, err := r.GetAccountByID(ctx, randomID)
	assert.Equal(t, entities.ErrAccountDoesNotExist, err)

	acc := entities.NewAccount(
		"Maria",
		testdata.GetValidCPF(),
		testdata.GetValidSecret(),
	)

	err = r.CreateAccount(ctx, &acc)
	assert.NoError(t, err)

	got, err := r.GetAccountByID(ctx, acc.ID)
	assert.NoError(t, err)
	assert.Equal(t, acc.ID, got.ID)
	assert.Equal(t, acc.Name, got.Name)
	assert.Equal(t, acc.CPF, got.CPF)
	assert.Equal(t, acc.Balance, got.Balance)
	assert.Equal(t, acc.CreatedAt, got.CreatedAt)
}
