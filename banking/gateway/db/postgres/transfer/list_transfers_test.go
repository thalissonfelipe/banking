package transfer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/db/postgres/account"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestRepository_ListTransfers(t *testing.T) {
	t.Parallel()

	db := dockertest.NewDB(t, t.Name())
	accRepository := account.NewRepository(db)
	r := NewRepository(db)
	ctx := context.Background()

	transfers, err := r.ListTransfers(ctx, vos.NewAccountID())
	require.NoError(t, err)
	assert.Len(t, transfers, 0)

	accOrigin, err := entity.NewAccount("origin", testdata.CPF().String(), testdata.Secret().String())
	require.NoError(t, err)

	accDest, err := entity.NewAccount("dest", testdata.CPF().String(), testdata.Secret().String())
	require.NoError(t, err)

	err = accRepository.CreateAccount(ctx, &accOrigin)
	require.NoError(t, err)

	err = accRepository.CreateAccount(ctx, &accDest)
	require.NoError(t, err)

	// It must return empty because the accounts have not yet carried out
	// any transactions
	transfers, err = r.ListTransfers(ctx, accOrigin.ID)
	require.NoError(t, err)
	assert.Len(t, transfers, 0)

	transfer, err := entity.NewTransfer(accOrigin.ID, accDest.ID, 50, 100)
	require.NoError(t, err)

	err = r.PerformTransfer(ctx, &transfer)
	require.NoError(t, err)
	assert.True(t, transfer.CreatedAt.Before(time.Now()))

	transfers, err = r.ListTransfers(ctx, accOrigin.ID)
	require.NoError(t, err)

	assert.Len(t, transfers, 1)
	assert.Equal(t, transfer.AccountOriginID, transfers[0].AccountOriginID)
	assert.Equal(t, transfer.AccountDestinationID, transfers[0].AccountDestinationID)
	assert.Equal(t, transfer.Amount, transfers[0].Amount)
}
