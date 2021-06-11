package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_GetAccounts(t *testing.T) {
	r := Repository{collection: collection}

	defer dropCollection(t, collection)

	accounts, err := r.GetAccounts(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, accounts)

	acc := createAccount(t, 0)
	assert.NoError(t, err)

	accounts, err = r.GetAccounts(context.Background())
	assert.NoError(t, err)
	assert.Len(t, accounts, 1)

	assert.Equal(t, acc.ID, accounts[0].ID)
	assert.Equal(t, acc.Name, accounts[0].Name)
	assert.Equal(t, acc.CPF, accounts[0].CPF)
	assert.Equal(t, acc.Balance, accounts[0].Balance)
}
