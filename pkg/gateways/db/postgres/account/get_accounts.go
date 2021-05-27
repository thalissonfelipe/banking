package account

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

const getAccountsQuery = `
select id, name, cpf, balance, created_at from accounts
`

func (r Repository) GetAccounts(ctx context.Context) ([]entities.Account, error) {
	rows, err := r.db.Query(ctx, getAccountsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]entities.Account, 0)

	for rows.Next() {
		var account entities.Account
		err := rows.Scan(
			&account.ID,
			&account.Name,
			&account.CPF,
			&account.Balance,
			&account.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
