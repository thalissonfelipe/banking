package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/tests/mocks"
)

func TestCreateAccount(t *testing.T) {
	ctx := context.Background()
	input := account.CreateAccountInput{
		Name:   "Pedro",
		CPF:    "648.446.967-93",
		Secret: "aZ1234Ds",
	}

	t.Run("should create an account successfully", func(t *testing.T) {
		repo := mocks.StubAccountRepository{}
		enc := mocks.StubHash{}
		usecase := NewAccountUseCase(&repo, enc)
		result, err := usecase.CreateAccount(ctx, input)

		assert.Nil(t, err)
		assert.Equal(t, input.Name, result.Name)
		assert.Equal(t, input.CPF, result.CPF)
		assert.NotEqual(t, input.Secret, result.Secret)
		assert.Len(t, repo.Accounts, 1)
	})

	t.Run("should return an error if repository fails to fetch or save", func(t *testing.T) {
		repo := mocks.StubAccountRepository{Err: errors.New("failed to save account")}
		enc := mocks.StubHash{}
		usecase := NewAccountUseCase(&repo, enc)
		result, err := usecase.CreateAccount(ctx, input)

		assert.Nil(t, result)
		assert.Equal(t, entities.ErrInternalError, err)
	})

	t.Run("should return an error if cpf already exists", func(t *testing.T) {
		acc := entities.NewAccount(input.Name, input.CPF, input.Secret)
		repo := mocks.StubAccountRepository{
			Accounts: []entities.Account{acc},
		}
		enc := mocks.StubHash{}
		usecase := NewAccountUseCase(&repo, enc)
		result, err := usecase.CreateAccount(ctx, input)

		assert.Nil(t, result)
		assert.Equal(t, entities.ErrAccountAlreadyExists, err)
	})

	t.Run("should return an error if hash secret fails", func(t *testing.T) {
		repo := mocks.StubAccountRepository{}
		enc := mocks.StubHash{Err: errors.New("could not hash secret")}
		usecase := NewAccountUseCase(&repo, enc)
		result, err := usecase.CreateAccount(ctx, input)

		assert.Nil(t, result)
		assert.Equal(t, entities.ErrInternalError, err)
	})

	t.Run("should return an error if cpf is not valid", func(t *testing.T) {
		input.CPF = "123.456.789-00"
		repo := mocks.StubAccountRepository{}
		enc := mocks.StubHash{}
		usecase := NewAccountUseCase(&repo, enc)
		result, err := usecase.CreateAccount(ctx, input)

		assert.Nil(t, result)
		assert.Equal(t, entities.ErrInvalidCPF, err)
	})

	t.Run("should return an error if secret is not valid", func(t *testing.T) {
		input.CPF = "648.446.967-93"
		input.Secret = "invalid_secret"
		repo := mocks.StubAccountRepository{}
		enc := mocks.StubHash{}
		usecase := NewAccountUseCase(&repo, enc)
		result, err := usecase.CreateAccount(ctx, input)

		assert.Nil(t, result)
		assert.Equal(t, entities.ErrInvalidSecret, err)
	})
}
