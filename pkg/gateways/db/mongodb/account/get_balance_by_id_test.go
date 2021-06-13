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
	r := NewRepository(db)

	defer dockertest.DropCollection(t, db.Collection("accounts"))

	balance, err := r.GetBalanceByID(context.Background(), vos.NewID())
	assert.Zero(t, balance)
	assert.ErrorIs(t, err, entities.ErrAccountDoesNotExist)

	acc := dockertest.CreateAccount(t, db.Collection("accounts"), 100)

	balance, err = r.GetBalanceByID(context.Background(), acc.ID)
	assert.NoError(t, err)
	assert.Equal(t, acc.Balance, balance)
}
