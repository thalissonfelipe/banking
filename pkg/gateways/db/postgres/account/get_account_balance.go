package account

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

const getBalanceQuery = `
select balance from accounts where id=$1
`

func (r Repository) GetBalanceByID(ctx context.Context, id vos.ID) (int, error) {
	var balance int

	err := r.db.QueryRow(ctx, getBalanceQuery, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, entities.ErrAccountDoesNotExist
		}
		return 0, err
	}

	return balance, nil
}
