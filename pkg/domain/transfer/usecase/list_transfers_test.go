package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestUsecase_ListTransfers(t *testing.T) {
	accountOriginID := vos.NewID()

	testCases := []struct {
		name        string
		repoSetup   func() *mocks.StubTransferRepository
		accountID   vos.ID
		expectedErr error
	}{
		{
			name: "should return a list of transfers",
			repoSetup: func() *mocks.StubTransferRepository {
				transfer := entities.NewTransfer(
					accountOriginID,
					vos.NewID(),
					100,
				)

				return &mocks.StubTransferRepository{
					Transfers: []entities.Transfer{transfer},
				}
			},
			accountID:   accountOriginID,
			expectedErr: nil,
		},
		{
			name: "should return an error if something went wrong on repository",
			repoSetup: func() *mocks.StubTransferRepository {
				return &mocks.StubTransferRepository{
					Err: testdata.ErrRepositoryFailsToFetch,
				}
			},
			accountID:   accountOriginID,
			expectedErr: entities.ErrInternalError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			usecase := NewTransferUsecase(tt.repoSetup(), nil)
			_, err := usecase.ListTransfers(ctx, tt.accountID)

			// TODO: add result validation
			assert.Equal(t, err, tt.expectedErr)
		})
	}
}
