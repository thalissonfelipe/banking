package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/tests/mocks"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestUsecase_ListAccounts(t *testing.T) {
	acc := entities.NewAccount("Piter", testdata.GetValidCPF(), testdata.GetValidSecret())

	testCases := []struct {
		name        string
		repoSetup   *mocks.AccountRepositoryMock
		expected    []entities.Account
		expectedErr error
	}{
		{
			name: "should return a list of accounts",
			repoSetup: &mocks.AccountRepositoryMock{
				Accounts: []entities.Account{acc},
			},
			expected:    []entities.Account{acc},
			expectedErr: nil,
		},
		{
			name:        "should return an error if something went wrong on repository",
			repoSetup:   &mocks.AccountRepositoryMock{Err: testdata.ErrRepositoryFailsToFetch},
			expected:    nil,
			expectedErr: entities.ErrInternalError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			usecase := NewAccountUsecase(tt.repoSetup, nil)
			result, err := usecase.ListAccounts(ctx)

			assert.Equal(t, tt.expected, result)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
