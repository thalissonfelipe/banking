package account

import (
	"context"
	"database/sql"
	"errors"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (r Repository) GetBalanceByID(ctx context.Context, id string) (int, error) {
	query := "SELECT balance FROM account WHERE id=$1"

	var balance int
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, entities.ErrAccountDoesNotExist
		}
		return 0, err
	}

	return balance, nil
}
