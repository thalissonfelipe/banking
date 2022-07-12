package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
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
	).Scan(&account.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return fmt.Errorf("db.QueryRow.Scan: %w", err)
		}

		if pgErr.Code == pgerrcode.UniqueViolation {
			return entities.ErrAccountAlreadyExists
		}

		return fmt.Errorf("db.QueryRow.Scan: %w", err)
	}

	return nil
}
