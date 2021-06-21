package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestRepository_GetAccountByCPF(t *testing.T) {
	db := pgDocker.DB
	r := NewRepository(db)
	ctx := context.Background()

	defer dockertest.TruncateTables(ctx, db)

	_, err := r.GetAccountByCPF(ctx, testdata.GetValidCPF())

	assert.Equal(t, entities.ErrAccountDoesNotExist, err)

	acc := entities.NewAccount(
		"Maria",
		testdata.GetValidCPF(),
		testdata.GetValidSecret(),
	)

	err = r.CreateAccount(ctx, &acc)
	assert.NoError(t, err)

	got, err := r.GetAccountByCPF(ctx, acc.CPF)
	assert.NoError(t, err)
	assert.Equal(t, acc.ID, got.ID)
	assert.Equal(t, acc.Name, got.Name)
	assert.Equal(t, acc.CPF, got.CPF)
	assert.Equal(t, acc.Balance, got.Balance)
	assert.Equal(t, acc.CreatedAt, got.CreatedAt)
}
