package account

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

const getAccountByIDQuery = `
select id, name, cpf, secret, balance, created_at
from accounts
where id=$1
`

func (r Repository) GetAccountByID(ctx context.Context, id vos.ID) (*entities.Account, error) {
	var account entities.Account

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
			return nil, entities.ErrAccountDoesNotExist
		}
		return nil, err
	}

	return &account, nil
}
