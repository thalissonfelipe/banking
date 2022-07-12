package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/tests/mocks"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestUsecase_GetBalanceByAccountID(t *testing.T) {
	acc := entities.NewAccount("Piter", testdata.GetValidCPF(), testdata.GetValidSecret())

	testCases := []struct {
		name        string
		repoSetup   *mocks.AccountRepositoryMock
		accountID   vos.AccountID
		expected    int
		expectedErr error
	}{
		{
			name: "should return a balance successfully",
			repoSetup: &mocks.AccountRepositoryMock{
				Accounts: []entities.Account{acc},
			},
			accountID:   acc.ID,
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
