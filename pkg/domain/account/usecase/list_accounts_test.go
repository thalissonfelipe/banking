package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
)

func TestUsecase_ListAccounts(t *testing.T) {
	acc := entities.NewAccount("Piter", vos.NewCPF("123.456.789-00"), vos.NewSecret("12345678"))
	testCases := []struct {
		name        string
		repoSetup   *mocks.StubAccountRepository
		expected    []entities.Account
		errExpected error
	}{
		{
			name: "should return a list of accounts",
			repoSetup: &mocks.StubAccountRepository{
				Accounts: []entities.Account{acc},
				Err:      nil,
			},
			expected:    []entities.Account{acc},
			errExpected: nil,
		},
		{
			name: "should return an error if something went wrong on repository",
			repoSetup: &mocks.StubAccountRepository{
				Accounts: nil,
				Err:      errors.New("failed to fetch accounts"),
			},
			expected:    nil,
			errExpected: entities.ErrInternalError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			usecase := NewAccountUsecase(tt.repoSetup, nil)
			result, err := usecase.ListAccounts(ctx)

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.errExpected, err)
		})
	}
}
