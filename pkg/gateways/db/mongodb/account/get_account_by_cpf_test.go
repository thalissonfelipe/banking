package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestRepostory_GetAccountByCPF(t *testing.T) {
	r := NewRepository(db)

	defer dockertest.DropCollection(t, db.Collection("accounts"))

	account, err := r.GetAccountByCPF(context.Background(), testdata.GetValidCPF())
	assert.Nil(t, account)
	assert.ErrorIs(t, err, entities.ErrAccountDoesNotExist)

	acc := dockertest.CreateAccount(t, db.Collection("accounts"), 100)

	account, err = r.GetAccountByCPF(context.Background(), acc.CPF)
	assert.NoError(t, err)

	assertAccountResponse(t, acc, *account, false)
}
