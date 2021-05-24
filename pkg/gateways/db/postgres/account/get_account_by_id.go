package account

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (r Repository) GetAccountByID(ctx context.Context, id vos.ID) (*entities.Account, error) {
	const query = `
		SELECT
			id,
			name,
			cpf,
			secret,
			balance,
			created_at
		FROM
			accounts
		WHERE
			id=$1`

	var account entities.Account
	err := r.db.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.Name,
		&account.CPF,
		&account.Secret,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		log.WithError(err).Error("unable to get account by id")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrAccountDoesNotExist
		}
		return nil, err
	}

	return &account, nil
}
