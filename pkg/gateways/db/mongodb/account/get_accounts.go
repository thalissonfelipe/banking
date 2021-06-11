package account

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (r Repository) GetAccounts(ctx context.Context) ([]entities.Account, error) {
	opts := options.Find().SetProjection(bson.M{"secret": 0})

	cur, err := r.db.Collection("accounts").Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("could not get accounts: %w", err)
	}

	accounts := make([]entities.Account, 0)

	err = cur.All(ctx, &accounts)
	if err != nil {
		return nil, fmt.Errorf("could not decode cursor: %w", err)
	}

	if err = cur.Err(); err != nil {
		return nil, fmt.Errorf("unexpected cursor error: %w", err)
	}

	return accounts, nil
}
