package account

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

const createAccountQuery = `
insert into accounts (id, name, cpf, secret, balance)
values ($1, $2, $3, $4, $5)
returning created_at
`

func (r Repository) CreateAccount(ctx context.Context, account *entities.Account) error {
	err := r.db.QueryRow(ctx, createAccountQuery,
		account.ID,
		account.Name,
		account.CPF.String(),
		account.Secret.String(),
		account.Balance,
	).Scan(
		&account.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
