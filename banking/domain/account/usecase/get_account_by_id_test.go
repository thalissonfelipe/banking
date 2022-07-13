package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountUsecase_GetAccountByID(t *testing.T) {
	acc, err := entities.NewAccount("name", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
	require.NoError(t, err)

	testCases := []struct {
		name    string
		repo    account.Repository
		want    entities.Account
		wantErr error
	}{
		{
			name: "should return an account successfully",
			repo: &RepositoryMock{
				GetAccountByIDFunc: func(context.Context, vos.AccountID) (entities.Account, error) {
					return acc, nil
				},
			},
			want:    acc,
			wantErr: nil,
		},
		{
			name: "should return an error if account does not exist",
			repo: &RepositoryMock{
				GetAccountByIDFunc: func(context.Context, vos.AccountID) (entities.Account, error) {
					return entities.Account{}, entities.ErrAccountNotFound
				},
			},
			want:    entities.Account{},
			wantErr: entities.ErrAccountNotFound,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			usecase := NewAccountUsecase(tt.repo, nil)

			account, err := usecase.GetAccountByID(context.Background(), vos.NewAccountID())
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, account)
		})
	}
}
