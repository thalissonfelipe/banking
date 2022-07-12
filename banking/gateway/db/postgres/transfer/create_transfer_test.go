package transfer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/gateway/db/postgres/account"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestTransferRepository_CreateTransfer(t *testing.T) {
	db := pgDocker.DB
	accRepository := account.NewRepository(db)
	r := NewRepository(db)
	ctx := context.Background()

	defer dockertest.TruncateTables(ctx, db)

	accOrigin := entities.NewAccount("origin", testdata.GetValidCPF(), testdata.GetValidSecret())
	accOrigin.Balance = 100
	accDestination := entities.NewAccount("destination", testdata.GetValidCPF(), testdata.GetValidSecret())

	err := accRepository.CreateAccount(ctx, &accOrigin)
	require.NoError(t, err)

	err = accRepository.CreateAccount(ctx, &accDestination)
	require.NoError(t, err)

	transfer := entities.NewTransfer(accOrigin.ID, accDestination.ID, 50)

	err = r.CreateTransfer(ctx, &transfer)
	require.NoError(t, err)
	assert.True(t, transfer.CreatedAt.Before(time.Now()))

	transfers, err := r.GetTransfers(ctx, accOrigin.ID)
	require.NoError(t, err)

	assert.Len(t, transfers, 1)
	assert.Equal(t, transfer.ID, transfers[0].ID)
	assert.Equal(t, transfer.AccountOriginID, transfers[0].AccountOriginID)
	assert.Equal(t, transfer.AccountDestinationID, transfers[0].AccountDestinationID)
	assert.Equal(t, transfer.Amount, transfers[0].Amount)

	originBalance, err := accRepository.GetBalanceByID(ctx, accOrigin.ID)
	require.NoError(t, err)
	assert.Equal(t, 50, originBalance)

	destinationBalance, err := accRepository.GetBalanceByID(ctx, accDestination.ID)
	require.NoError(t, err)
	assert.Equal(t, 50, destinationBalance)
}
