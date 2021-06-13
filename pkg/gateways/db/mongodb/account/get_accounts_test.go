package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
)

func TestRepository_GetAccounts(t *testing.T) {
	r := NewRepository(db)

	defer dockertest.DropCollection(t, db.Collection("accounts"))

	accounts, err := r.GetAccounts(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, accounts)

	acc := dockertest.CreateAccount(t, db.Collection("accounts"), 0)
	assert.NoError(t, err)

	accounts, err = r.GetAccounts(context.Background())
	assert.NoError(t, err)
	assert.Len(t, accounts, 1)

	assertAccountResponse(t, acc, accounts[0], true)
}
