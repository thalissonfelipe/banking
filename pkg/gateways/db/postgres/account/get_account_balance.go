package account

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (r Repository) GetBalanceByID(ctx context.Context, id vos.ID) (int, error) {
	query := "SELECT balance FROM account WHERE id=$1"

	var balance int
	err := r.db.QueryRow(ctx, query, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, entities.ErrAccountDoesNotExist
		}
		return 0, err
	}

	return balance, nil
}
