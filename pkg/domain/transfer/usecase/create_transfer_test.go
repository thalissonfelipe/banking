package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
)

func TestCreateTransfer(t *testing.T) {
	accOrigin := entities.NewAccount("Pedro", vos.NewCPF("123.456.789-00"), vos.NewSecret("12345678"))
	accDest := entities.NewAccount("Maria", vos.NewCPF("123.456.789-01"), vos.NewSecret("12345678"))

	testCases := []struct {
		name       string
		repoSetup  *mocks.StubTransferRepository
		accUsecase func() *mocks.StubAccountUseCase
		input      func() transfer.CreateTransferInput
		expected   error
	}{
		{
			name:      "should perform a transfer successfully",
			repoSetup: &mocks.StubTransferRepository{},
			accUsecase: func() *mocks.StubAccountUseCase {
				accOriginWithBalance := accOrigin
				accOriginWithBalance.Balance = 100
				return &mocks.StubAccountUseCase{
					Accounts: []entities.Account{accOriginWithBalance, accDest},
				}
			},
			input: func() transfer.CreateTransferInput {
				return transfer.NewTransferInput(accOrigin.ID, accDest.ID, 100)
			},
			expected: nil,
		},
		{
			name:      "should return an error if accOrigin does not have sufficient funds",
			repoSetup: &mocks.StubTransferRepository{},
			accUsecase: func() *mocks.StubAccountUseCase {
				return &mocks.StubAccountUseCase{
					Accounts: []entities.Account{accOrigin, accDest},
				}
			},
			input: func() transfer.CreateTransferInput {
				return transfer.NewTransferInput(accOrigin.ID, accDest.ID, 100)
			},
			expected: entities.ErrInsufficientFunds,
		},
		{
			name: "should return an error if transfer repository fails",
			repoSetup: &mocks.StubTransferRepository{
				Err: errors.New("failed to transfer amount"),
			},
			accUsecase: func() *mocks.StubAccountUseCase {
				accOriginWithBalance := accOrigin
				accOriginWithBalance.Balance = 100
				return &mocks.StubAccountUseCase{
					Accounts: []entities.Account{accOriginWithBalance, accDest},
				}
			},
			input: func() transfer.CreateTransferInput {
				return transfer.NewTransferInput(accOrigin.ID, accDest.ID, 100)
			},
			expected: entities.ErrInternalError,
		},
		{
			name:      "should return an error if accUsecase fails",
			repoSetup: &mocks.StubTransferRepository{},
			accUsecase: func() *mocks.StubAccountUseCase {
				accOriginWithBalance := accOrigin
				accOriginWithBalance.Balance = 100
				return &mocks.StubAccountUseCase{
					Accounts: []entities.Account{accOriginWithBalance, accDest},
					Err:      errors.New("failed to get account origin balance"),
				}
			},
			input: func() transfer.CreateTransferInput {
				return transfer.NewTransferInput(accOrigin.ID, accDest.ID, 100)
			},
			expected: entities.ErrInternalError,
		},
		{
			name:      "should return an error if account origin does not exist",
			repoSetup: &mocks.StubTransferRepository{},
			accUsecase: func() *mocks.StubAccountUseCase {
				return &mocks.StubAccountUseCase{
					Accounts: []entities.Account{accDest},
				}
			},
			input: func() transfer.CreateTransferInput {
				return transfer.NewTransferInput(accOrigin.ID, accDest.ID, 100)
			},
			expected: entities.ErrAccountDoesNotExist,
		},
		{
			name:      "should return an error if account destination does not exist",
			repoSetup: &mocks.StubTransferRepository{},
			accUsecase: func() *mocks.StubAccountUseCase {
				accOriginWithBalance := accOrigin
				accOriginWithBalance.Balance = 100
				return &mocks.StubAccountUseCase{
					Accounts: []entities.Account{accOriginWithBalance},
				}
			},
			input: func() transfer.CreateTransferInput {
				return transfer.NewTransferInput(accOrigin.ID, accDest.ID, 100)
			},
			expected: entities.ErrAccountDestinationDoesNotExist,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			usecase := NewTransfer(tt.repoSetup, tt.accUsecase())
			err := usecase.CreateTransfer(ctx, tt.input())

			assert.Equal(t, tt.expected, err)
		})
	}
}
