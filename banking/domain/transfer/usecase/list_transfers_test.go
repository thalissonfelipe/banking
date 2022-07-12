package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/tests/mocks"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestUsecase_ListTransfers(t *testing.T) {
	accountOriginID := vos.NewAccountID()

	testCases := []struct {
		name      string
		repoSetup func() *mocks.TransferRepositoryMock
		accountID vos.AccountID
		wantErr   bool
	}{
		{
			name: "should return a list of transfers",
			repoSetup: func() *mocks.TransferRepositoryMock {
				transfer := entities.NewTransfer(
					accountOriginID,
					vos.NewAccountID(),
					100,
				)

				return &mocks.TransferRepositoryMock{
					Transfers: []entities.Transfer{transfer},
				}
			},
			accountID: accountOriginID,
			wantErr:   false,
		},
		{
			name: "should return an error if something went wrong on repository",
			repoSetup: func() *mocks.TransferRepositoryMock {
				return &mocks.TransferRepositoryMock{
					Err: testdata.ErrRepositoryFailsToFetch,
				}
			},
			accountID: accountOriginID,
			wantErr:   true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			usecase := NewTransferUsecase(tt.repoSetup(), nil)
			_, err := usecase.ListTransfers(ctx, tt.accountID)

			// TODO: add result validation
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
