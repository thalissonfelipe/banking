package transfer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateways/db/postgres/account"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestRepository_GetTransfers(t *testing.T) {
	db := pgDocker.DB
	accRepository := account.NewRepository(db)
	r := NewRepository(db)
	ctx := context.Background()

	randomID := vos.NewAccountID()

	transfers, err := r.GetTransfers(ctx, randomID)
	assert.NoError(t, err)
	assert.Len(t, transfers, 0)

	acc1 := entities.NewAccount(
		"Maria",
		testdata.GetValidCPF(),
		testdata.GetValidSecret(),
	)
	acc1.Balance = 100

	acc2 := entities.NewAccount(
		"Pedro",
		testdata.GetValidCPF(),
		testdata.GetValidSecret(),
	)

	err = accRepository.CreateAccount(ctx, &acc1)
	assert.NoError(t, err)

	err = accRepository.CreateAccount(ctx, &acc2)
	assert.NoError(t, err)

	// It must return empty because the accounts have not yet carried out
	// any transactions
	transfers, err = r.GetTransfers(ctx, acc1.ID)
	assert.NoError(t, err)
	assert.Len(t, transfers, 0)

	transfer := entities.NewTransfer(acc1.ID, acc2.ID, 50)

	err = r.CreateTransfer(ctx, &transfer)
	assert.NoError(t, err)

	transfers, err = r.GetTransfers(ctx, acc1.ID)

	assert.NoError(t, err)
	assert.Len(t, transfers, 1)
	assert.Equal(t, transfer.AccountOriginID, transfers[0].AccountOriginID)
	assert.Equal(t, transfer.AccountDestinationID, transfers[0].AccountDestinationID)
	assert.Equal(t, transfer.Amount, transfers[0].Amount)
}
