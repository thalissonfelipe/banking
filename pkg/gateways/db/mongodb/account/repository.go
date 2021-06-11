package account

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
)

var _ account.Repository = (*Repository)(nil)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) Repository {
	return Repository{db: db}
}
