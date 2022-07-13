package transfer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/db/postgres/account"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestRepository_GetTransfers(t *testing.T) {
	db := pgDocker.DB
	accRepository := account.NewRepository(db)
	r := NewRepository(db)
	ctx := context.Background()

	defer dockertest.TruncateTables(ctx, db)

	transfers, err := r.GetTransfers(ctx, vos.NewAccountID())
	require.NoError(t, err)
	assert.Len(t, transfers, 0)

	accOrigin, err := entities.NewAccount("origin", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
	require.NoError(t, err)

	accDest, err := entities.NewAccount("dest", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
	require.NoError(t, err)

	err = accRepository.CreateAccount(ctx, &accOrigin)
	require.NoError(t, err)

	err = accRepository.CreateAccount(ctx, &accDest)
	require.NoError(t, err)

	// It must return empty because the accounts have not yet carried out
	// any transactions
	transfers, err = r.GetTransfers(ctx, accOrigin.ID)
	require.NoError(t, err)
	assert.Len(t, transfers, 0)

	transfer := entities.NewTransfer(accOrigin.ID, accDest.ID, 50)

	err = r.CreateTransfer(ctx, &transfer)
	require.NoError(t, err)
	assert.True(t, transfer.CreatedAt.Before(time.Now()))

	transfers, err = r.GetTransfers(ctx, accOrigin.ID)
	require.NoError(t, err)

	assert.Len(t, transfers, 1)
	assert.Equal(t, transfer.AccountOriginID, transfers[0].AccountOriginID)
	assert.Equal(t, transfer.AccountDestinationID, transfers[0].AccountDestinationID)
	assert.Equal(t, transfer.Amount, transfers[0].Amount)
}
