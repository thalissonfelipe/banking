package transfer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
)

func TestRepository_GetTransfers(t *testing.T) {
	r := NewRepository(db)

	defer dockertest.DropCollection(t, db.Collection("transfers"))
	defer dockertest.DropCollection(t, db.Collection("accounts"))

	transfers, err := r.GetTransfers(context.Background(), vos.NewAccountID())
	assert.NoError(t, err)
	assert.Len(t, transfers, 0)

	acc1 := dockertest.CreateAccount(t, db.Collection("accounts"), 100)
	acc2 := dockertest.CreateAccount(t, db.Collection("accounts"), 100)

	transfer := entities.NewTransfer(acc1.ID, acc2.ID, 100)
	transfer.CreatedAt = time.Now()

	err = r.CreateTransfer(context.Background(), &transfer)
	assert.NoError(t, err)

	// still should return an empty slice because the account origin id does not exist.
	transfers, err = r.GetTransfers(context.Background(), vos.NewAccountID())
	assert.NoError(t, err)
	assert.Len(t, transfers, 0)

	transfers, err = r.GetTransfers(context.Background(), transfer.AccountOriginID)
	assert.NoError(t, err)
	assert.Len(t, transfers, 1)

	assert.Equal(t, transfer.ID, transfers[0].ID)
	assert.Equal(t, transfer.AccountOriginID, transfers[0].AccountOriginID)
	assert.Equal(t, transfer.AccountDestinationID, transfers[0].AccountDestinationID)
	assert.Equal(t, transfer.Amount, transfers[0].Amount)
	assert.Equal(t, transfer.CreatedAt.Unix(), transfers[0].CreatedAt.Unix())
}
