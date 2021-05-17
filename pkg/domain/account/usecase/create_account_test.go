package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/account"
)

func TestCreateAccount(t *testing.T) {
	ctx := context.Background()

	t.Run("should create an account", func(t *testing.T) {
		input := account.CreateAccountInput{
			Name:   "Pedro",
			CPF:    "123.456.789-00",
			Secret: "12345678",
		}
		repo := StubRepository{}
		usecase := NewAccountUseCase(&repo)
		result, err := usecase.CreateAccount(ctx, input)

		assert.Nil(t, err)
		assert.Equal(t, input.Name, result.Name)
		assert.Equal(t, input.CPF, result.CPF)
		assert.Len(t, repo.accounts, 1)
	})

	t.Run("should return an error if repository fails to save", func(t *testing.T) {
		input := account.CreateAccountInput{
			Name:   "Pedro",
			CPF:    "123.456.789-00",
			Secret: "12345678",
		}
		repo := StubRepository{err: errors.New("failed to save account")}
		usecase := NewAccountUseCase(&repo)
		result, err := usecase.CreateAccount(ctx, input)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Len(t, repo.accounts, 0)
	})
}
