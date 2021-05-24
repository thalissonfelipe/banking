package account

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (r Repository) GetBalanceByID(ctx context.Context, id vos.ID) (int, error) {
	query := `
		SELECT
			balance
		FROM accounts
		WHERE id=$1
	`

	var balance int
	err := r.db.QueryRow(ctx, query, id).Scan(&balance)
	if err != nil {
		log.WithError(err).Error("unable to get balance by id")
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, entities.ErrAccountDoesNotExist
		}
		return 0, err
	}

	return balance, nil
}
