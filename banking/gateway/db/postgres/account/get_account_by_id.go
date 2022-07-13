package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

const getAccountByIDQuery = `
select id, name, cpf, secret, balance, created_at
from accounts
where id=$1
`

func (r Repository) GetAccountByID(ctx context.Context, id vos.AccountID) (entity.Account, error) {
	var account entity.Account

	err := r.db.QueryRow(ctx, getAccountByIDQuery, id).Scan(
		&account.ID,
		&account.Name,
		&account.CPF,
		&account.Secret,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Account{}, entity.ErrAccountNotFound
		}

		return entity.Account{}, fmt.Errorf("db.QueryRow.Scan: %w", err)
	}

	return account, nil
}
