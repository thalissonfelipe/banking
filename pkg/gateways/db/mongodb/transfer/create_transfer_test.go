package transfer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/gateways/db/mongodb/account"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
)

func TestRepository_CreateTransfer(t *testing.T) {
	r := NewRepository(db)
	ctx := context.Background()

	defer dockertest.DropCollection(t, db.Collection("transfers"))
	defer dockertest.DropCollection(t, db.Collection("accounts"))

	acc1 := dockertest.CreateAccount(t, db.Collection("accounts"), 100)
	acc2 := dockertest.CreateAccount(t, db.Collection("accounts"), 100)

	transfer := entities.NewTransfer(acc1.ID, acc2.ID, 100)
	transfer.CreatedAt = time.Now()

	err := r.CreateTransfer(ctx, &transfer)
	assert.NoError(t, err)

	transfers, err := r.GetTransfers(ctx, acc1.ID)
	assert.NoError(t, err)

	assert.Equal(t, transfer.ID, transfers[0].ID)
	assert.Equal(t, transfer.AccountOriginID, transfers[0].AccountOriginID)
	assert.Equal(t, transfer.AccountDestinationID, transfers[0].AccountDestinationID)
	assert.Equal(t, transfer.Amount, transfers[0].Amount)
	assert.Equal(t, transfer.CreatedAt.Unix(), transfers[0].CreatedAt.Unix())

	accountRepository := account.NewRepository(db)

	accounts, err := accountRepository.GetAccounts(context.Background())
	assert.NoError(t, err)

	if accounts[0].ID == acc1.ID {
		assert.Equal(t, 0, accounts[0].Balance)
		assert.Equal(t, 200, accounts[1].Balance)
	} else {
		assert.Equal(t, 0, accounts[1].Balance)
		assert.Equal(t, 200, accounts[0].Balance)
	}
}
