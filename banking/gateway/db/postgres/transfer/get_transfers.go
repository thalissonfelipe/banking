package transfer

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

const getTransfersQuery = `
select id, account_origin_id, account_destination_id, amount, created_at
from transfers
where account_origin_id=$1
`

func (r Repository) GetTransfers(ctx context.Context, id vos.AccountID) ([]entities.Transfer, error) {
	rows, err := r.db.Query(ctx, getTransfersQuery, id)
	if err != nil {
		return nil, fmt.Errorf("unexpected error occurred on get transfers query: %w", err)
	}
	defer rows.Close()

	transfers := make([]entities.Transfer, 0)

	for rows.Next() {
		var account entities.Transfer

		err = rows.Scan(
			&account.ID,
			&account.AccountOriginID,
			&account.AccountDestinationID,
			&account.Amount,
			&account.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("unexpected error occurred while scanning rows: %w", err)
		}

		transfers = append(transfers, account)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("unexpected error occurred while scanning rows: %w", err)
	}

	return transfers, nil
}
