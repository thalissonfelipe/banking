package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
)

func TestCreateAccount(t *testing.T) {
	validInput := account.NewCreateAccountInput("Pedro", vos.NewCPF("648.446.967-93"), "aZ1234Ds")
	invalidSecretInput := account.NewCreateAccountInput("João", vos.NewCPF("648.446.967-93"), "12345678")

	testCases := []struct {
		name        string
		repoSetup   func() *mocks.StubAccountRepository
		encSetup    *mocks.StubHash
		input       account.CreateAccountInput
		errExpected error
	}{
		{
			name: "should create an account successfully",
			repoSetup: func() *mocks.StubAccountRepository {
				return &mocks.StubAccountRepository{}
			},
			encSetup:    &mocks.StubHash{},
			input:       validInput,
			errExpected: nil,
		},
		{
			name: "should return an error if repository fails to fetch or save",
			repoSetup: func() *mocks.StubAccountRepository {
				return &mocks.StubAccountRepository{
					Err: errors.New("failed to save account"),
				}
			},
			encSetup:    &mocks.StubHash{},
			input:       validInput,
			errExpected: entities.ErrInternalError,
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
			errExpected: entities.ErrAccountAlreadyExists,
		},
		{
			name: "should return an error if hash secret fails",
			repoSetup: func() *mocks.StubAccountRepository {
				return &mocks.StubAccountRepository{}
			},
			input:       validInput,
			encSetup:    &mocks.StubHash{Err: errors.New("could not hash secret")},
			errExpected: entities.ErrInternalError,
		},
		{
			name: "should return an error if secret is not valid",
			repoSetup: func() *mocks.StubAccountRepository {
				return &mocks.StubAccountRepository{}
			},
			input:       invalidSecretInput,
			encSetup:    &mocks.StubHash{},
			errExpected: entities.ErrInvalidSecret,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			usecase := NewAccountUseCase(tt.repoSetup(), tt.encSetup)
			_, err := usecase.CreateAccount(ctx, tt.input)

			// TODO: add result validation
			assert.Equal(t, tt.errExpected, err)
		})
	}
}
