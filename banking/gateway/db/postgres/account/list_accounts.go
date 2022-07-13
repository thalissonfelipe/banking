package account

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
)

const listAccountsQuery = `
select id, name, cpf, balance, created_at from accounts
`

func (r Repository) ListAccounts(ctx context.Context) ([]entity.Account, error) {
	rows, err := r.db.Query(ctx, listAccountsQuery)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()

	accounts := make([]entity.Account, 0)

	for rows.Next() {
		var account entity.Account

		err = rows.Scan(
			&account.ID,
			&account.Name,
			&account.CPF,
			&account.Balance,
			&account.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Scan: %w", err)
	}

	return accounts, nil
}
