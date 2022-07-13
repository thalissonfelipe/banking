package account

import (
	"github.com/jackc/pgx/v4"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
)

var _ entity.AccountRepository = (*Repository)(nil)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}
