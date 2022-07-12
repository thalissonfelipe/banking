package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/tests/mocks"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
)

func TestUsecase_CreateAccount(t *testing.T) {
	validInput := account.NewCreateAccountInput(
		"Pedro",
		testdata.GetValidCPF(),
		testdata.GetValidSecret(),
	)

	testCases := []struct {
		name        string
		repoSetup   func() *mocks.AccountRepositoryMock
		encSetup    *mocks.HashMock
		input       account.CreateAccountInput
		expectedErr error
	}{
		{
			name: "should create an account successfully",
			repoSetup: func() *mocks.AccountRepositoryMock {
				return &mocks.AccountRepositoryMock{}
			},
			encSetup:    &mocks.HashMock{},
			input:       validInput,
			expectedErr: nil,
		},
		{
			name: "should return an error if account already exists",
			repoSetup: func() *mocks.AccountRepositoryMock {
				acc := entities.NewAccount(validInput.Name, validInput.CPF, validInput.Secret)

				return &mocks.AccountRepositoryMock{
					Accounts: []entities.Account{acc},
				}
			},
			input:       validInput,
			encSetup:    &mocks.HashMock{},
			expectedErr: entities.ErrAccountAlreadyExists,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			usecase := NewAccountUsecase(tt.repoSetup(), tt.encSetup)
			_, err := usecase.CreateAccount(ctx, tt.input)

			// TODO: add result validation
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
