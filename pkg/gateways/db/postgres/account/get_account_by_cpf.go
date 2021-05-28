package account

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

const getAccountByCPFQuery = `
select id, name, cpf, secret, balance, created_at
from accounts
where cpf=$1
`

func (r Repository) GetAccountByCPF(ctx context.Context, cpf vos.CPF) (*entities.Account, error) {
	var account entities.Account

	err := r.db.QueryRow(ctx, getAccountByCPFQuery, cpf).Scan(
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
