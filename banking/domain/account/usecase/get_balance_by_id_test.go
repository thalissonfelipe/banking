package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func TestAccountUsecase_GetAccountBalanceByID(t *testing.T) {
	wantBalance := 100

	testCases := []struct {
		name    string
		repo    account.Repository
		want    int
		wantErr error
	}{
		{
			name: "should return an account balance successfully",
			repo: &RepositoryMock{
				GetAccountBalanceByIDFunc: func(context.Context, vos.AccountID) (int, error) {
					return wantBalance, nil
				},
			},
			want:    wantBalance,
			wantErr: nil,
		},
		{
			name: "should return an error if account does not exist",
			repo: &RepositoryMock{
				GetAccountBalanceByIDFunc: func(context.Context, vos.AccountID) (int, error) {
					return 0, entities.ErrAccountNotFound
				},
			},
			want:    0,
			wantErr: entities.ErrAccountNotFound,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			usecase := NewAccountUsecase(tt.repo, nil)

			balance, err := usecase.GetAccountBalanceByID(context.Background(), vos.NewAccountID())
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, balance)
		})
	}
}
