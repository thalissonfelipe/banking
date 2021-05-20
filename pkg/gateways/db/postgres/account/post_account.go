package account

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (r Repository) PostAccount(ctx context.Context, account entities.Account) error {
	const query = `
		INSERT INTO account (
			id,
			name,
			cpf,
			secret,
			balance,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`

	_, err := r.db.Exec(ctx, query,
		account.ID,
		account.Name,
		account.CPF.String(),
		account.Secret.String(),
		account.Balance,
		account.CreatedAt,
	)

	return err
}
