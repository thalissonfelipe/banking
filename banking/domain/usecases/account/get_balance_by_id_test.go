package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func TestAccountUsecase_GetAccountBalanceByID(t *testing.T) {
	t.Parallel()

	wantBalance := 100

	testCases := []struct {
		name    string
		repo    entity.AccountRepository
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
					return 0, entity.ErrAccountNotFound
				},
			},
			want:    0,
			wantErr: entity.ErrAccountNotFound,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			usecase := NewAccountUsecase(tt.repo, nil)

			balance, err := usecase.GetAccountBalanceByID(context.Background(), vos.NewAccountID())
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, balance)
		})
	}
}
