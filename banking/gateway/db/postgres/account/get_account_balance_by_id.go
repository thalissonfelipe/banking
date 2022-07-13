package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

const getBalanceQuery = `
select balance from accounts where id=$1
`

func (r Repository) GetAccountBalanceByID(ctx context.Context, id vos.AccountID) (int, error) {
	var balance int

	err := r.db.QueryRow(ctx, getBalanceQuery, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, entity.ErrAccountNotFound
		}

		return 0, fmt.Errorf("db.QueryRow.Scan: %w", err)
	}

	return balance, nil
}
