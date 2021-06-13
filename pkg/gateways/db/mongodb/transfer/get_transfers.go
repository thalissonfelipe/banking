package transfer

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

type transferAdapter struct {
	ID                   vos.ID    `bson:"id"`
	AccountOriginID      vos.ID    `bson:"account_origin_id"`
	AccountDestinationID vos.ID    `bson:"account_destination_id"`
	Amount               int       `bson:"amount"`
	CreatedAt            time.Time `bson:"created_at"`
}

func (r Repository) GetTransfers(ctx context.Context, id vos.ID) ([]entities.Transfer, error) {
	cur, err := r.db.Collection("transfers").Find(ctx, bson.M{"account_origin_id": id})
	if err != nil {
		return nil, fmt.Errorf("could not get transfers: %w", err)
	}

	transfersBSON := make([]transferAdapter, 0)

	err = cur.All(ctx, &transfersBSON)
	if err != nil {
		return nil, fmt.Errorf("could not decode cursor: %w", err)
	}

	if err = cur.Err(); err != nil {
		return nil, fmt.Errorf("unexpected cursor error: %w", err)
	}

	transfers := make([]entities.Transfer, 0)

	for _, t := range transfersBSON {
		transfers = append(transfers, entities.Transfer{
			ID:                   t.ID,
			AccountOriginID:      t.AccountOriginID,
			AccountDestinationID: t.AccountDestinationID,
			Amount:               t.Amount,
			CreatedAt:            t.CreatedAt,
		})
	}

	return transfers, nil
}
