package transfer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestRepository_GetTransfers(t *testing.T) {
	r := NewRepository(collection)

	defer dockertest.DropCollection(t, collection)

	transfers, err := r.GetTransfers(context.Background(), vos.NewID())
	assert.NoError(t, err)
	assert.Len(t, transfers, 0)

	acc1 := entities.NewAccount("Felipe", testdata.GetValidCPF(), testdata.GetValidSecret())
	acc2 := entities.NewAccount("Sousa", testdata.GetValidCPF(), testdata.GetValidSecret())

	transfer := entities.NewTransfer(acc1.ID, acc2.ID, 100)
	transfer.CreatedAt = time.Now()

	insertedID, err := collection.InsertOne(context.Background(), bson.D{
		primitive.E{Key: "id", Value: transfer.ID},
		primitive.E{Key: "account_origin_id", Value: transfer.AccountOriginID},
		primitive.E{Key: "account_destination_id", Value: transfer.AccountDestinationID},
		primitive.E{Key: "amount", Value: transfer.Amount},
		primitive.E{Key: "created_at", Value: transfer.CreatedAt},
	})
	assert.NoError(t, err)
	assert.NotNil(t, insertedID)

	// still should return an empty slice because the account origin id does not exist.
	transfers, err = r.GetTransfers(context.Background(), vos.NewID())
	assert.NoError(t, err)
	assert.Len(t, transfers, 0)

	transfers, err = r.GetTransfers(context.Background(), transfer.AccountOriginID)
	assert.NoError(t, err)
	assert.Len(t, transfers, 1)
}
