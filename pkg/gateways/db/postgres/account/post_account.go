package account

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (r Repository) PostAccount(ctx context.Context, account *entities.Account) error {
	const query = `
		INSERT INTO accounts (
			id,
			name,
			cpf,
			secret,
			balance
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING created_at
	`

	err := r.db.QueryRow(ctx, query,
		account.ID,
		account.Name,
		account.CPF.String(),
		account.Secret.String(),
		account.Balance,
	).Scan(
		&account.CreatedAt,
	)
	if err != nil {
		log.WithError(err).Error("unable to create a new account")
		return err
	}

	return nil
}
