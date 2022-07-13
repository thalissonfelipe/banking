package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/transfer"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func TestTransferUsecase_ListTransfers(t *testing.T) {
	tr, err := entities.NewTransfer(vos.NewAccountID(), vos.NewAccountID(), 100, 100)
	require.NoError(t, err)

	transfers := []entities.Transfer{tr}

	testCases := []struct {
		name    string
		repo    transfer.Repository
		wantErr error
	}{
		{
			name: "should return a list of transfers successfully",
			repo: &RepositoryMock{
				GetTransfersFunc: func(context.Context, vos.AccountID) ([]entities.Transfer, error) {
					return transfers, nil
				},
			},
			wantErr: nil,
		},
		{
			name: "should return an error if repo fails to get transfers",
			repo: &RepositoryMock{
				GetTransfersFunc: func(context.Context, vos.AccountID) ([]entities.Transfer, error) {
					return nil, assert.AnError
				},
			},
			wantErr: assert.AnError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			usecase := NewTransferUsecase(tt.repo, nil)

			_, err := usecase.ListTransfers(context.Background(), vos.NewAccountID())
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
