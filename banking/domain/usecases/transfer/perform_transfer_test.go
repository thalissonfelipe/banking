package transfer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestTransferUsecase_PerformTransfer(t *testing.T) {
	accOrigin, err := entity.NewAccount("origin", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
	require.NoError(t, err)

	accDest, err := entity.NewAccount("dest", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
	require.NoError(t, err)

	testCases := []struct {
		name       string
		repo       entity.TransferRepository
		accUsecase usecases.Account
		input      usecases.PerformTransferInput
		wantErr    error
	}{
		{
			name: "should perform a transfer successfully",
			repo: &RepositoryMock{
				PerformTransferFunc: func(context.Context, *entity.Transfer) error {
					return nil
				},
			},
			accUsecase: &UsecaseMock{
				GetAccountByIDFunc: func(_ context.Context, id vos.AccountID) (entity.Account, error) {
					if id == accOrigin.ID {
						return accOrigin, nil
					}

					return accDest, nil
				},
			},
			input:   usecases.NewPerformTransferInput(accOrigin.ID, accDest.ID, 100),
			wantErr: nil,
		},
		{
			name: "should return an error if account origin does not have sufficient funds",
			repo: &RepositoryMock{},
			accUsecase: &UsecaseMock{
				GetAccountByIDFunc: func(_ context.Context, id vos.AccountID) (entity.Account, error) {
					if id == accOrigin.ID {
						return accOrigin, nil
					}

					return accDest, nil
				},
			},
			input:   usecases.NewPerformTransferInput(accOrigin.ID, accDest.ID, accOrigin.Balance+1),
			wantErr: entity.ErrInsufficientFunds,
		},
		{
			name: "should return an error if account origin does not exist",
			repo: &RepositoryMock{},
			accUsecase: &UsecaseMock{
				GetAccountByIDFunc: func(context.Context, vos.AccountID) (entity.Account, error) {
					return entity.Account{}, entity.ErrAccountNotFound
				},
			},
			input:   usecases.NewPerformTransferInput(accOrigin.ID, accDest.ID, accOrigin.Balance),
			wantErr: entity.ErrAccountNotFound,
		},
		{
			name: "should return an error if account destination does not exist",
			repo: &RepositoryMock{},
			accUsecase: &UsecaseMock{
				GetAccountByIDFunc: func(_ context.Context, id vos.AccountID) (entity.Account, error) {
					if id == accOrigin.ID {
						return accOrigin, nil
					}

					return entity.Account{}, entity.ErrAccountNotFound
				},
			},
			input:   usecases.NewPerformTransferInput(accOrigin.ID, accDest.ID, accOrigin.Balance),
			wantErr: entity.ErrAccountDestinationNotFound,
		},
		{
			name: "should return an error if account usecase fails",
			repo: &RepositoryMock{},
			accUsecase: &UsecaseMock{
				GetAccountByIDFunc: func(context.Context, vos.AccountID) (entity.Account, error) {
					return entity.Account{}, assert.AnError
				},
			},
			input:   usecases.NewPerformTransferInput(accOrigin.ID, accDest.ID, accOrigin.Balance),
			wantErr: assert.AnError,
		},
		{
			name: "should return an error if repo fails to create a transfer",
			repo: &RepositoryMock{
				PerformTransferFunc: func(context.Context, *entity.Transfer) error {
					return assert.AnError
				},
			},
			accUsecase: &UsecaseMock{
				GetAccountByIDFunc: func(_ context.Context, id vos.AccountID) (entity.Account, error) {
					if id == accOrigin.ID {
						return accOrigin, nil
					}

					return accDest, nil
				},
			},
			input:   usecases.NewPerformTransferInput(accOrigin.ID, accDest.ID, accOrigin.Balance),
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
