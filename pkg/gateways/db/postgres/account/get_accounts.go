package account

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (r Repository) GetAccounts(ctx context.Context) ([]entities.Account, error) {
	const query = `
		SELECT
			id,
			name,
			cpf,
			balance,
			created_at
		FROM accounts
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		log.WithError(err).Error("unable to get accounts")
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
			log.WithError(err).Error("unable to scan account when iterating over rows")
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		log.WithError(err).Error("unexpected error ocurred while reading rows")
		return nil, err
	}

	return accounts, nil
}
