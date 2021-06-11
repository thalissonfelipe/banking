package account

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
)

var _ account.Repository = (*Repository)(nil)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) Repository {
	return Repository{collection: collection}
}
