package transfer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/gateway/db/postgres/account"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestTransferRepository_PerformTransfer(t *testing.T) {
	db := pgDocker.DB
	accRepository := account.NewRepository(db)
	r := NewRepository(db)
	ctx := context.Background()

	defer dockertest.TruncateTables(ctx, db)

	accOrigin, err := entity.NewAccount("origin", testdata.CPF().String(), testdata.Secret().String())
	require.NoError(t, err)

	accDest, err := entity.NewAccount("dest", testdata.CPF().String(), testdata.Secret().String())
	require.NoError(t, err)

	err = accRepository.CreateAccount(ctx, &accOrigin)
	require.NoError(t, err)

	err = accRepository.CreateAccount(ctx, &accDest)
	require.NoError(t, err)

	transfer, err := entity.NewTransfer(accOrigin.ID, accDest.ID, 50, 100)
	require.NoError(t, err)

	err = r.PerformTransfer(ctx, &transfer)
	require.NoError(t, err)
	assert.True(t, transfer.CreatedAt.Before(time.Now()))

	transfers, err := r.ListTransfers(ctx, accOrigin.ID)
	require.NoError(t, err)

	assert.Len(t, transfers, 1)
	assert.Equal(t, transfer.ID, transfers[0].ID)
	assert.Equal(t, transfer.AccountOriginID, transfers[0].AccountOriginID)
	assert.Equal(t, transfer.AccountDestinationID, transfers[0].AccountDestinationID)
	assert.Equal(t, transfer.Amount, transfers[0].Amount)

	originBalance, err := accRepository.GetAccountBalanceByID(ctx, accOrigin.ID)
	require.NoError(t, err)
	assert.Equal(t, 50, originBalance)

	destinationBalance, err := accRepository.GetAccountBalanceByID(ctx, accDest.ID)
	require.NoError(t, err)
	assert.Equal(t, 150, destinationBalance)
}
