package transfer

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (r Repository) GetTransfers(ctx context.Context, id vos.ID) ([]entities.Transfer, error) {
	cur, err := r.db.Collection("transfers").Find(ctx, bson.M{"account_origin_id": id})
	if err != nil {
		return nil, fmt.Errorf("could not get transfers: %w", err)
	}

	transfers := make([]entities.Transfer, 0)

	err = cur.All(ctx, &transfers)
	if err != nil {
		return nil, fmt.Errorf("could not decode cursor: %w", err)
	}

	if err = cur.Err(); err != nil {
		return nil, fmt.Errorf("unexpected cursor error: %w", err)
	}

	return transfers, nil
}
