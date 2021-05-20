package account

import (
	"context"
	"database/sql"
	"errors"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (r Repository) GetAccountByID(ctx context.Context, id string) (*entities.Account, error) {
	const query = `
		SELECT
			id,
			name,
			cpf,
			secret,
			balance,
			created_at
		FROM
			account
		WHERE
			id=$1`

	var account entities.Account
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.Name,
		&account.CPF,
		&account.Secret,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrAccountDoesNotExist
		}
		return nil, err
	}

	return &account, nil
}
