package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func TestRepository_CreateAccount(t *testing.T) {
	db := pgDocker.DB
	r := NewRepository(db)
	ctx := context.Background()

	acc := entities.NewAccount(
		"Maria",
		vos.NewCPF("123.456.789-00"),
		vos.NewSecret("12345678"),
	)

	assert.Empty(t, acc.CreatedAt)

	err := r.CreateAccount(ctx, &acc)
	assert.NoError(t, err)
	assert.NotEmpty(t, acc.CreatedAt)

	// Should return an error if cpf already exists
	err = r.CreateAccount(ctx, &acc)
	assert.Equal(t, err, entities.ErrAccountAlreadyExists)

	account, err := r.GetAccountByID(ctx, acc.ID)

	assert.NoError(t, err)
	assert.Equal(t, acc.Name, account.Name)
	assert.Equal(t, acc.CPF, account.CPF)
	assert.Equal(t, acc.Balance, account.Balance)
}
