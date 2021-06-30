package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
)

func TestRepostory_GetAccountByID(t *testing.T) {
	r := NewRepository(db)

	defer dockertest.DropCollection(t, db.Collection("accounts"))

	account, err := r.GetAccountByID(context.Background(), vos.NewAccountID())
	assert.Nil(t, account)
	assert.ErrorIs(t, err, entities.ErrAccountDoesNotExist)

	acc := dockertest.CreateAccount(t, db.Collection("accounts"), 100)

	account, err = r.GetAccountByID(context.Background(), acc.ID)
	assert.NoError(t, err)

	assertAccountResponse(t, acc, *account, false)
}
