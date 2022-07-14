package transfer

import (
	"context"
	"fmt"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

const _listTransfersQuery = `
select
	id,
	account_origin_id,
	account_destination_id,
	amount,
	created_at
from transfers
where account_origin_id=$1;`

func (r Repository) ListTransfers(ctx context.Context, id vos.AccountID) ([]entity.Transfer, error) {
	rows, err := r.db.Query(ctx, _listTransfersQuery, id)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()

	transfers := make([]entity.Transfer, 0)

	for rows.Next() {
		var account entity.Transfer

		if err = rows.Scan(
			&account.ID,
			&account.AccountOriginID,
			&account.AccountDestinationID,
			&account.Amount,
			&account.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		transfers = append(transfers, account)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Scan: %w", err)
	}

	return transfers, nil
}
