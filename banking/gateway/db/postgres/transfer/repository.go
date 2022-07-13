package transfer

import (
	"github.com/jackc/pgx/v4"

	"github.com/thalissonfelipe/banking/banking/domain/transfer"
)

var _ transfer.Repository = (*Repository)(nil)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}
