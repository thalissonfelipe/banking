package transfer

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (r Repository) UpdateBalance(ctx context.Context, transfer entities.Transfer) error {
	// First experience with rollback.
	// Tutorial: https://www.sohamkamani.com/golang/sql-transactions/
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

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

	tx.Commit(ctx)

	return nil
}

func (r Repository) updateBalance(ctx context.Context, tx pgx.Tx, balance int, id vos.ID) error {
	const query = `UPDATE account SET balance=balance+$1 WHERE id=$2`

	_, err := tx.Exec(ctx, query, balance, id)
	return err
}

func (r Repository) saveTransfer(ctx context.Context, tx pgx.Tx, transfer entities.Transfer) error {
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

	_, err := tx.Exec(ctx, query,
		transfer.ID,
		transfer.AccountOriginID,
		transfer.AccountDestinationID,
		transfer.Amount,
		transfer.CreatedAt,
	)

	return err
}
