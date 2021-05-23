package transfer

import (
	"context"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (r Repository) GetTransfers(ctx context.Context, id vos.ID) ([]entities.Transfer, error) {
	const query = `
		SELECT
			id,
			account_origin_id,
			account_destination_id,
			amount,
			created_at
		FROM transfer
		WHERE account_origin_id=$1
	`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transfers := make([]entities.Transfer, 0)

	for rows.Next() {
		var account entities.Transfer
		err := rows.Scan(
			&account.ID,
			&account.AccountOriginID,
			&account.AccountDestinationID,
			&account.Amount,
			&account.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transfers, nil
}
