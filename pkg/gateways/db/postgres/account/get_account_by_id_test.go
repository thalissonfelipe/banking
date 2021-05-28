package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
)

func TestRepository_GetAccountByID(t *testing.T) {
	db := pgDocker.DB
	r := NewRepository(db)
	ctx := context.Background()

	defer dockertest.TruncateTables(ctx, db)

	randomID := vos.NewID()
	a, err := r.GetAccountByID(ctx, randomID)

	assert.Equal(t, entities.ErrAccountDoesNotExist, err)
	assert.Nil(t, a)

	acc := entities.NewAccount(
		"Maria",
		vos.NewCPF("123.456.789-00"),
		vos.NewSecret("12345678"),
	)

	err = r.CreateAccount(ctx, &acc)
	assert.NoError(t, err)

	a, err = r.GetAccountByID(ctx, acc.ID)
	assert.NoError(t, err)
	assert.Equal(t, acc.ID, a.ID)
	assert.Equal(t, acc.Name, a.Name)
	assert.Equal(t, acc.CPF, a.CPF)
	assert.Equal(t, acc.Balance, a.Balance)
	assert.Equal(t, acc.CreatedAt, a.CreatedAt)
}
