package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
)

func TestListTransfers(t *testing.T) {
	accountOriginID := entities.NewAccountID()

	testCases := []struct {
		name        string
		repoSetup   func() *mocks.StubTransferRepository
		accountId   string
		errExpected error
	}{
		{
			name: "should return a list of transfers",
			repoSetup: func() *mocks.StubTransferRepository {
				transfer := entities.NewTransfer(
					accountOriginID,
					entities.NewAccountID(),
					100,
				)
				return &mocks.StubTransferRepository{
					Transfers: []entities.Transfer{transfer},
				}
			},
			accountId:   accountOriginID,
			errExpected: nil,
		},
		{
			name: "should return an error if something went wrong on repository",
			repoSetup: func() *mocks.StubTransferRepository {
				return &mocks.StubTransferRepository{
					Err: errors.New("failed to fetch transfers"),
				}
			},
			accountId:   accountOriginID,
			errExpected: entities.ErrInternalError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			usecase := NewTransfer(tt.repoSetup(), nil)
			_, err := usecase.ListTransfers(ctx, tt.accountId)

			// TODO: add result validation
			assert.Equal(t, tt.errExpected, err)
		})
	}
}
