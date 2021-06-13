package account

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (r Repository) GetBalanceByID(ctx context.Context, accountID vos.ID) (int, error) {
	opts := options.FindOne().SetProjection(bson.M{"balance": 1})

	var account accountAdpater

	err := r.db.Collection("accounts").FindOne(ctx, bson.M{"id": accountID}, opts).Decode(&account)
	if err == nil {
		return account.Balance, nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return 0, entities.ErrAccountDoesNotExist
	}

	return 0, fmt.Errorf("could not get account balance: %w", err)
}
