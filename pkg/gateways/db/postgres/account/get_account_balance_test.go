package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
)

func TestRepository_GetBalanceByID(t *testing.T) {
	db := pgDocker.DB
	r := NewRepository(db)
	ctx := context.Background()

	defer dockertest.TruncateTables(ctx, db)

	randomID := vos.NewID()
	balance, err := r.GetBalanceByID(ctx, randomID)

	assert.Equal(t, entities.ErrAccountDoesNotExist, err)
	assert.Equal(t, 0, balance)

	acc := entities.NewAccount(
		"Maria",
		vos.NewCPF("123.456.789-00"),
		vos.NewSecret("12345678"),
	)

	err = r.PostAccount(ctx, &acc)
	assert.NoError(t, err)

	balance, err = r.GetBalanceByID(ctx, acc.ID)

	assert.NoError(t, err)
	assert.Equal(t, acc.Balance, balance)
}
