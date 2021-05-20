package transfer

import (
	"context"
	"database/sql"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (r Repository) UpdateBalance(ctx context.Context, transfer entities.Transfer) error {
	// First experience with rollback.
	// Tutorial: https://www.sohamkamani.com/golang/sql-transactions/
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = r.updateBalance(ctx, tx, -transfer.Amount, transfer.AccountOriginID)
	if err != nil {
		return err
	}

	err = r.updateBalance(ctx, tx, transfer.Amount, transfer.AccountDestinationID)
	if err != nil {
		return err
	}

	err = r.saveTransfer(ctx, tx, transfer)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (r Repository) updateBalance(ctx context.Context, tx *sql.Tx, balance int, id string) error {
	const query = `UPDATE account SET balance=balance+$1 WHERE id=$2`

	_, err := tx.ExecContext(ctx, query, balance, id)
	return err
}

func (r Repository) saveTransfer(ctx context.Context, tx *sql.Tx, transfer entities.Transfer) error {
	const query = `
		INSERT INTO transfer (
			id,
			account_origin_id,
			account_destination_id,
			amount,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5
		)
	`

	_, err := tx.ExecContext(ctx, query,
		transfer.ID,
		transfer.AccountOriginID,
		transfer.AccountDestinationID,
		transfer.Amount,
		transfer.CreatedAt,
	)

	return err
}
