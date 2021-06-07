package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestUsecase_CreateTransfer(t *testing.T) {
	accOrigin := entities.NewAccount("Pedro", testdata.GetValidCPF(), testdata.GetValidSecret())
	accDest := entities.NewAccount("Maria", testdata.GetValidCPF(), testdata.GetValidSecret())

	testCases := []struct {
		name        string
		repoSetup   *mocks.StubTransferRepository
		accUsecase  func() *mocks.StubAccountUsecase
		input       func() transfer.CreateTransferInput
		expectedErr error
	}{
		{
			name:      "should perform a transfer successfully",
			repoSetup: &mocks.StubTransferRepository{},
			accUsecase: func() *mocks.StubAccountUsecase {
				accOriginWithBalance := accOrigin
				accOriginWithBalance.Balance = 100

				return &mocks.StubAccountUsecase{
					Accounts: []entities.Account{accOriginWithBalance, accDest},
				}
			},
			input: func() transfer.CreateTransferInput {
				return transfer.NewTransferInput(accOrigin.ID, accDest.ID, 100)
			},
			expectedErr: nil,
		},
		{
			name:      "should return an error if accOrigin does not have sufficient funds",
			repoSetup: &mocks.StubTransferRepository{},
			accUsecase: func() *mocks.StubAccountUsecase {
				return &mocks.StubAccountUsecase{
					Accounts: []entities.Account{accOrigin, accDest},
				}
			},
			input: func() transfer.CreateTransferInput {
				return transfer.NewTransferInput(accOrigin.ID, accDest.ID, 100)
			},
			expectedErr: entities.ErrInsufficientFunds,
		},
		{
			name:      "should return an error if transfer repository fails",
			repoSetup: &mocks.StubTransferRepository{Err: testdata.ErrRepositoryFailsToSave},
			accUsecase: func() *mocks.StubAccountUsecase {
				accOriginWithBalance := accOrigin
				accOriginWithBalance.Balance = 100

				return &mocks.StubAccountUsecase{
					Accounts: []entities.Account{accOriginWithBalance, accDest},
				}
			},
			input: func() transfer.CreateTransferInput {
				return transfer.NewTransferInput(accOrigin.ID, accDest.ID, 100)
			},
			expectedErr: entities.ErrInternalError,
		},
		{
			name:      "should return an error if accUsecase fails",
			repoSetup: &mocks.StubTransferRepository{},
			accUsecase: func() *mocks.StubAccountUsecase {
				accOriginWithBalance := accOrigin
				accOriginWithBalance.Balance = 100

				return &mocks.StubAccountUsecase{
					Accounts: []entities.Account{accOriginWithBalance, accDest},
					Err:      testdata.ErrUsecaseFails,
				}
			},
			input: func() transfer.CreateTransferInput {
				return transfer.NewTransferInput(accOrigin.ID, accDest.ID, 100)
			},
			expectedErr: entities.ErrInternalError,
		},
		{
			name:      "should return an error if account origin does not exist",
			repoSetup: &mocks.StubTransferRepository{},
			accUsecase: func() *mocks.StubAccountUsecase {
				return &mocks.StubAccountUsecase{
					Accounts: []entities.Account{accDest},
				}
			},
			input: func() transfer.CreateTransferInput {
				return transfer.NewTransferInput(accOrigin.ID, accDest.ID, 100)
			},
			expectedErr: entities.ErrAccountDoesNotExist,
		},
		{
			name:      "should return an error if account destination does not exist",
			repoSetup: &mocks.StubTransferRepository{},
			accUsecase: func() *mocks.StubAccountUsecase {
				accOriginWithBalance := accOrigin
				accOriginWithBalance.Balance = 100

				return &mocks.StubAccountUsecase{
					Accounts: []entities.Account{accOriginWithBalance},
				}
			},
			input: func() transfer.CreateTransferInput {
				return transfer.NewTransferInput(accOrigin.ID, accDest.ID, 100)
			},
			expectedErr: entities.ErrAccountDestinationDoesNotExist,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			usecase := NewTransferUsecase(tt.repoSetup, tt.accUsecase())
			err := usecase.CreateTransfer(ctx, tt.input())

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
