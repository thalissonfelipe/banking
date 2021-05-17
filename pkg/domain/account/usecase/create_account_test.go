package usecase

import (
	"context"
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
		usecase := NewAccountUseCase(repo)
		result, err := usecase.CreateAccount(ctx, input)

		assert.Nil(t, err)
		assert.Equal(t, input.Name, result.Name)
		assert.Equal(t, input.CPF, result.CPF)
	})
}
