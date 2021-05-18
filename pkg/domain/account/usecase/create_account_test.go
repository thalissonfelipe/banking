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
		CPF:    "123.456.789-00",
		Secret: "12345678",
	}

	t.Run("should create an account", func(t *testing.T) {
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
		assert.NotNil(t, err)
		assert.Len(t, repo.Accounts, 0)
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
		assert.NotNil(t, err)
	})
}
