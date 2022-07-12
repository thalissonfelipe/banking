package transfer

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func (r Repository) CreateTransfer(ctx context.Context, transfer *entities.Transfer) error {
	// First experience with rollback.
	// Tutorial: https://www.sohamkamani.com/golang/sql-transactions/
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err = r.updateBalance(ctx, tx, -transfer.Amount, transfer.AccountOriginID); err != nil {
		return err
	}

	if err = r.updateBalance(ctx, tx, transfer.Amount, transfer.AccountDestinationID); err != nil {
		return err
	}

	if err = r.saveTransfer(ctx, tx, transfer); err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}

const updateBalanceQuery = `
update accounts set balance=balance+$1 where id=$2
`

func (r Repository) updateBalance(ctx context.Context, tx pgx.Tx, balance int, id vos.AccountID) error {
	_, err := tx.Exec(ctx, updateBalanceQuery, balance, id)
	if err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}

	return nil
}

const insertTransferQuery = `
insert into transfers (id, account_origin_id, account_destination_id, amount)
values ($1, $2, $3, $4)
returning created_at
`

func (r Repository) saveTransfer(ctx context.Context, tx pgx.Tx, transfer *entities.Transfer) error {
	err := tx.QueryRow(ctx, insertTransferQuery,
		transfer.ID,
		transfer.AccountOriginID,
		transfer.AccountDestinationID,
		transfer.Amount,
	).Scan(
		&transfer.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("tx.QueryRow.Scan: %w", err)
	}

	return nil
}
