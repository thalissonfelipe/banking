package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/transfer"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestTransferUsecase_PerformTransfer(t *testing.T) {
	accOrigin, err := entities.NewAccount("origin", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
	require.NoError(t, err)

	accDest, err := entities.NewAccount("dest", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
	require.NoError(t, err)

	testCases := []struct {
		name       string
		repo       transfer.Repository
		accUsecase account.Usecase
		input      transfer.PerformTransferInput
		wantErr    error
	}{
		{
			name: "should perform a transfer successfully",
			repo: &RepositoryMock{
				PerformTransferFunc: func(context.Context, *entities.Transfer) error {
					return nil
				},
			},
			accUsecase: &UsecaseMock{
				GetAccountByIDFunc: func(_ context.Context, id vos.AccountID) (entities.Account, error) {
					if id == accOrigin.ID {
						return accOrigin, nil
					}

					return accDest, nil
				},
			},
			input:   transfer.NewPerformTransferInput(accOrigin.ID, accDest.ID, 100),
			wantErr: nil,
		},
		{
			name: "should return an error if account origin does not have sufficient funds",
			repo: &RepositoryMock{},
			accUsecase: &UsecaseMock{
				GetAccountByIDFunc: func(_ context.Context, id vos.AccountID) (entities.Account, error) {
					if id == accOrigin.ID {
						return accOrigin, nil
					}

					return accDest, nil
				},
			},
			input:   transfer.NewPerformTransferInput(accOrigin.ID, accDest.ID, accOrigin.Balance+1),
			wantErr: entities.ErrInsufficientFunds,
		},
		{
			name: "should return an error if account origin does not exist",
			repo: &RepositoryMock{},
			accUsecase: &UsecaseMock{
				GetAccountByIDFunc: func(context.Context, vos.AccountID) (entities.Account, error) {
					return entities.Account{}, entities.ErrAccountNotFound
				},
			},
			input:   transfer.NewPerformTransferInput(accOrigin.ID, accDest.ID, accOrigin.Balance),
			wantErr: entities.ErrAccountNotFound,
		},
		{
			name: "should return an error if account destination does not exist",
			repo: &RepositoryMock{},
			accUsecase: &UsecaseMock{
				GetAccountByIDFunc: func(_ context.Context, id vos.AccountID) (entities.Account, error) {
					if id == accOrigin.ID {
						return accOrigin, nil
					}

					return entities.Account{}, entities.ErrAccountNotFound
				},
			},
			input:   transfer.NewPerformTransferInput(accOrigin.ID, accDest.ID, accOrigin.Balance),
			wantErr: entities.ErrAccountDestinationNotFound,
		},
		{
			name: "should return an error if account usecase fails",
			repo: &RepositoryMock{},
			accUsecase: &UsecaseMock{
				GetAccountByIDFunc: func(context.Context, vos.AccountID) (entities.Account, error) {
					return entities.Account{}, assert.AnError
				},
			},
			input:   transfer.NewPerformTransferInput(accOrigin.ID, accDest.ID, accOrigin.Balance),
			wantErr: assert.AnError,
		},
		{
			name: "should return an error if repo fails to create a transfer",
			repo: &RepositoryMock{
				PerformTransferFunc: func(context.Context, *entities.Transfer) error {
					return assert.AnError
				},
			},
			accUsecase: &UsecaseMock{
				GetAccountByIDFunc: func(_ context.Context, id vos.AccountID) (entities.Account, error) {
					if id == accOrigin.ID {
						return accOrigin, nil
					}

					return accDest, nil
				},
			},
			input:   transfer.NewPerformTransferInput(accOrigin.ID, accDest.ID, accOrigin.Balance),
			wantErr: assert.AnError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			usecase := NewTransferUsecase(tt.repo, tt.accUsecase)

			err := usecase.PerformTransfer(context.Background(), tt.input)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
