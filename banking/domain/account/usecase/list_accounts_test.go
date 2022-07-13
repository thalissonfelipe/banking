package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountUsecase_ListAccounts(t *testing.T) {
	acc, err := entity.NewAccount("name", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
	require.NoError(t, err)

	accounts := []entity.Account{acc}

	testCases := []struct {
		name    string
		repo    account.Repository
		want    []entity.Account
		wantErr error
	}{
		{
			name: "should return a list of accounts",
			repo: &RepositoryMock{
				ListAccountsFunc: func(context.Context) ([]entity.Account, error) {
					return accounts, nil
				},
			},
			want:    accounts,
			wantErr: nil,
		},
		{
			name: "should return an error if repo fails to get accounts",
			repo: &RepositoryMock{
				ListAccountsFunc: func(context.Context) ([]entity.Account, error) {
					return nil, assert.AnError
				},
			},
			want:    nil,
			wantErr: assert.AnError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			usecase := NewAccountUsecase(tt.repo, nil)

			accounts, err := usecase.ListAccounts(context.Background())
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, accounts)
		})
	}
}
