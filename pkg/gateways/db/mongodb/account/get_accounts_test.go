package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestRepository_GetAccounts(t *testing.T) {
	r := Repository{collection: collection}

	accounts, err := r.GetAccounts(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, accounts)

	acc := entities.NewAccount("Felipe", testdata.GetValidCPF(), testdata.GetValidSecret())

	_, err = collection.InsertOne(context.Background(), bson.D{
		primitive.E{Key: "id", Value: acc.ID},
		primitive.E{Key: "name", Value: acc.Name},
		primitive.E{Key: "cpf", Value: acc.CPF.String()},
		primitive.E{Key: "secret", Value: acc.Secret.String()},
		primitive.E{Key: "balance", Value: acc.Balance},
	})
	assert.NoError(t, err)

	accounts, err = r.GetAccounts(context.Background())
	assert.NoError(t, err)
	assert.Len(t, accounts, 1)

	assert.Equal(t, acc.ID, accounts[0].ID)
	assert.Equal(t, acc.Name, accounts[0].Name)
	assert.Equal(t, acc.CPF, accounts[0].CPF)
	assert.Equal(t, acc.Balance, accounts[0].Balance)
}
