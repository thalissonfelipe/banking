package account

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
)

var _ entity.AccountRepository = (*Repository)(nil)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}
