package transfer

import (
	"context"

	log "github.com/sirupsen/logrus"

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
		FROM transfers
		WHERE account_origin_id=$1
	`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		log.WithError(err).Error("unable to get transfers")
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
			log.WithError(err).Error("unable to scan transfer when iterating over rows")
			return nil, err
		}
		transfers = append(transfers, account)
	}

	if err := rows.Err(); err != nil {
		log.WithError(err).Error("unexpected error ocurred while reading rows")
		return nil, err
	}

	return transfers, nil
}
