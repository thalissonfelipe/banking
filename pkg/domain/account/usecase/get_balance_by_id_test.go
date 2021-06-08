package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestUsecase_GetBalanceByAccountID(t *testing.T) {
	accBalanceDefault := entities.NewAccount("Piter", testdata.GetValidCPF(), testdata.GetValidSecret())
	accBalance100 := entities.NewAccount("Piter", testdata.GetValidCPF(), testdata.GetValidSecret())
	accBalance100.Balance = 100

	testCases := []struct {
		name        string
		repoSetup   *mocks.AccountRepositoryMock
		accountID   vos.AccountID
		expected    int
		expectedErr error
	}{
		{
			name: "should return a default balance",
			repoSetup: &mocks.AccountRepositoryMock{
				Accounts: []entities.Account{accBalanceDefault},
			},
			accountID:   accBalanceDefault.ID,
			expected:    entities.DefaultBalance,
			expectedErr: nil,
		},
		{
			name:        "should return an error if account does not exist",
			repoSetup:   &mocks.AccountRepositoryMock{},
			accountID:   vos.NewAccountID(),
			expected:    entities.DefaultBalance,
			expectedErr: entities.ErrAccountDoesNotExist,
		},
		{
			name: "should return correct balance when balance is not default",
			repoSetup: &mocks.AccountRepositoryMock{
				Accounts: []entities.Account{accBalance100},
			},
			accountID:   accBalance100.ID,
			expected:    100,
			expectedErr: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			usecase := NewAccountUsecase(tt.repoSetup, nil)
			result, err := usecase.GetAccountBalanceByID(ctx, tt.accountID)

			assert.Equal(t, tt.expected, result)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
