package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestUsecase_CreateAccount(t *testing.T) {
	validInput := account.NewCreateAccountInput(
		"Pedro",
		testdata.GetValidCPF(),
		testdata.GetValidSecret(),
	)

	testCases := []struct {
		name        string
		repoSetup   func() *mocks.StubAccountRepository
		encSetup    *mocks.StubHash
		input       account.CreateAccountInput
		expectedErr error
	}{
		{
			name: "should create an account successfully",
			repoSetup: func() *mocks.StubAccountRepository {
				return &mocks.StubAccountRepository{}
			},
			encSetup:    &mocks.StubHash{},
			input:       validInput,
			expectedErr: nil,
		},
		{
			name: "should return an error if repository fails to fetch account",
			repoSetup: func() *mocks.StubAccountRepository {
				return &mocks.StubAccountRepository{Err: testdata.ErrRepositoryFailsToFetch}
			},
			encSetup:    &mocks.StubHash{},
			input:       validInput,
			expectedErr: entities.ErrInternalError,
		},
		{
			name: "should return an error if cpf already exists",
			repoSetup: func() *mocks.StubAccountRepository {
				acc := entities.NewAccount(validInput.Name, validInput.CPF, validInput.Secret)

				return &mocks.StubAccountRepository{
					Accounts: []entities.Account{acc},
				}
			},
			input:       validInput,
			encSetup:    &mocks.StubHash{},
			expectedErr: entities.ErrAccountAlreadyExists,
		},
		{
			name: "should return an error if hash secret fails",
			repoSetup: func() *mocks.StubAccountRepository {
				return &mocks.StubAccountRepository{}
			},
			input:       validInput,
			encSetup:    &mocks.StubHash{Err: testdata.ErrHashFails},
			expectedErr: entities.ErrInternalError,
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
