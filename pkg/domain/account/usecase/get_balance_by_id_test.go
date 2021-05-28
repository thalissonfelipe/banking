package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
)

func TestUsecase_GetBalanceByAccountID(t *testing.T) {
	accBalanceDefault := entities.NewAccount("Piter", vos.NewCPF("123.456.789-00"), vos.NewSecret("12345678"))
	accBalance100 := entities.NewAccount("Piter", vos.NewCPF("123.456.789-00"), vos.NewSecret("12345678"))
	accBalance100.Balance = 100

	testCases := []struct {
		name        string
		repoSetup   *mocks.StubAccountRepository
		accountId   vos.ID
		expected    int
		errExpected error
	}{
		{
			name: "should return a default balance",
			repoSetup: &mocks.StubAccountRepository{
				Accounts: []entities.Account{accBalanceDefault},
				Err:      nil,
			},
			accountId:   accBalanceDefault.ID,
			expected:    entities.DefaultBalance,
			errExpected: nil,
		},
		{
			name: "should return an error if account does not exist",
			repoSetup: &mocks.StubAccountRepository{
				Accounts: nil,
				Err:      nil,
			},
			accountId:   vos.NewID(),
			expected:    entities.DefaultBalance,
			errExpected: entities.ErrAccountDoesNotExist,
		},
		{
			name: "should return correct balance when balance is not default",
			repoSetup: &mocks.StubAccountRepository{
				Accounts: []entities.Account{accBalance100},
				Err:      nil,
			},
			accountId:   accBalance100.ID,
			expected:    100,
			errExpected: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			usecase := NewAccountUsecase(tt.repoSetup, nil)
			result, err := usecase.GetAccountBalanceByID(ctx, tt.accountId)

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.errExpected, err)
		})
	}
}
