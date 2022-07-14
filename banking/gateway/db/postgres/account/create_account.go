package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/thalissonfelipe/banking/banking/domain/entity"
)

const _createAccountQuery = `
insert into accounts (id, name, cpf, secret, balance)
values ($1, $2, $3, $4, $5)
returning created_at;`

func (r Repository) CreateAccount(ctx context.Context, account *entity.Account) error {
	if err := r.db.QueryRow(ctx, _createAccountQuery,
		account.ID,
		account.Name,
		account.CPF.String(),
		account.Secret.String(),
		account.Balance,
	).Scan(&account.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return fmt.Errorf("db.QueryRow.Scan: %w", err)
		}

		if pgErr.Code == pgerrcode.UniqueViolation {
			return entity.ErrAccountAlreadyExists
		}

		return fmt.Errorf("db.QueryRow.Scan: %w", err)
	}

	return nil
}
