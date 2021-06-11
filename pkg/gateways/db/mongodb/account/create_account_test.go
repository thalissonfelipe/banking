package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestRepository_CreateAccount(t *testing.T) {
	r := Repository{collection: collection}

	defer dropCollection(t, collection)

	account := entities.NewAccount("Felipe", testdata.GetValidCPF(), testdata.GetValidSecret())

	err := r.CreateAccount(context.Background(), &account)
	assert.NoError(t, err)
	assert.NotEmpty(t, account.CreatedAt)

	err = r.CreateAccount(context.Background(), &account)
	assert.ErrorIs(t, err, entities.ErrAccountAlreadyExists)
}
