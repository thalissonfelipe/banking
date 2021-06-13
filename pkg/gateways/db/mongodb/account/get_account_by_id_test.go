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

	account, err := r.GetAccountByID(context.Background(), vos.NewID())
	assert.Nil(t, account)
	assert.ErrorIs(t, err, entities.ErrAccountDoesNotExist)

	acc := dockertest.CreateAccount(t, db.Collection("accounts"), 100)

	account, err = r.GetAccountByID(context.Background(), acc.ID)
	assert.NoError(t, err)

	assert.Equal(t, acc.ID, account.ID)
	assert.Equal(t, acc.Name, account.Name)
	assert.Equal(t, acc.CPF, account.CPF)
	assert.Equal(t, acc.Balance, account.Balance)
	assert.Equal(t, acc.Secret, account.Secret)
	assert.Equal(t, acc.CreatedAt.Unix(), account.CreatedAt.Unix())
}
