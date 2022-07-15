package transfer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func TestTransferUsecase_ListTransfers(t *testing.T) {
	t.Parallel()

	tr, err := entity.NewTransfer(vos.NewAccountID(), vos.NewAccountID(), 100, 100)
	require.NoError(t, err)

	transfers := []entity.Transfer{tr}

	testCases := []struct {
		name    string
		repo    entity.TransferRepository
		wantErr error
	}{
		{
			name: "should return a list of transfers successfully",
			repo: &RepositoryMock{
				ListTransfersFunc: func(context.Context, vos.AccountID) ([]entity.Transfer, error) {
					return transfers, nil
				},
			},
			wantErr: nil,
		},
		{
			name: "should return an error if repo fails to get transfers",
			repo: &RepositoryMock{
				ListTransfersFunc: func(context.Context, vos.AccountID) ([]entity.Transfer, error) {
					return nil, assert.AnError
				},
			},
			wantErr: assert.AnError,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			usecase := NewTransferUsecase(tt.repo, nil)

			_, err := usecase.ListTransfers(context.Background(), vos.NewAccountID())
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
