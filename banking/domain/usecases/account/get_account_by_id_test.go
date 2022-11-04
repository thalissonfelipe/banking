package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestAccountUsecase_GetAccountByID(t *testing.T) {
	t.Parallel()

	acc, err := entity.NewAccount("name", testdata.CPF().String(), testdata.Secret().String())
	require.NoError(t, err)

	testCases := []struct {
		name    string
		repo    entity.AccountRepository
		want    entity.Account
		wantErr error
	}{
		{
			name: "should return an account successfully",
			repo: &RepositoryMock{
				GetAccountByIDFunc: func(context.Context, vos.AccountID) (entity.Account, error) {
					return acc, nil
				},
			},
			want:    acc,
			wantErr: nil,
		},
		{
			name: "should return an error if account does not exist",
			repo: &RepositoryMock{
				GetAccountByIDFunc: func(context.Context, vos.AccountID) (entity.Account, error) {
					return entity.Account{}, entity.ErrAccountNotFound
				},
			},
			want:    entity.Account{},
			wantErr: entity.ErrAccountNotFound,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			usecase := NewUsecase(tt.repo, nil)

			account, err := usecase.GetAccountByID(context.Background(), vos.NewAccountID())
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, account)
		})
	}
}
