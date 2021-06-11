package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func TestRepository_GetBalanceByID(t *testing.T) {
	r := Repository{collection: collection}

	defer dropCollection(t, collection)

	balance, err := r.GetBalanceByID(context.Background(), vos.NewID())
	assert.Zero(t, balance)
	assert.ErrorIs(t, err, entities.ErrAccountDoesNotExist)

	acc := createAccount(t, 100)

	balance, err = r.GetBalanceByID(context.Background(), acc.ID)
	assert.NoError(t, err)
	assert.Equal(t, acc.Balance, balance)
}
