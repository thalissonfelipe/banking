package transfer

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (r Repository) CreateTransfer(ctx context.Context, transfer *entities.Transfer) error {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	transfersCollection := r.db.Collection("transfers")
	accountCollection := r.db.Collection("accounts")

	callback := func(sc mongo.SessionContext) (interface{}, error) {
		err := r.updateBalance(sc, accountCollection, transfer.AccountOriginID, -transfer.Amount)
		if err != nil {
			return nil, err
		}

		err = r.updateBalance(sc, accountCollection, transfer.AccountDestinationID, transfer.Amount)
		if err != nil {
			return nil, err
		}

		err = r.createTransfer(sc, transfersCollection, transfer)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	session, err := r.db.Client().StartSession()
	if err != nil {
		return fmt.Errorf("could not start a new session: %w", err)
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, callback, txnOpts)
	if err != nil {
		return fmt.Errorf("could not run session.WithTransaction: %w", err)
	}

	return nil
}

func (r Repository) updateBalance(ctx context.Context, coll *mongo.Collection, id vos.AccountID, amount int) error {
	filter := bson.M{"id": id}
	update := bson.D{primitive.E{Key: "$inc", Value: bson.D{primitive.E{Key: "balance", Value: amount}}}}

	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("could not update balance: %w", err)
	}

	return nil
}

func (r Repository) createTransfer(ctx context.Context, coll *mongo.Collection, transfer *entities.Transfer) error {
	transfer.CreatedAt = time.Now()

	_, err := coll.InsertOne(ctx, bson.D{
		primitive.E{Key: "id", Value: transfer.ID},
		primitive.E{Key: "account_origin_id", Value: transfer.AccountOriginID},
		primitive.E{Key: "account_destination_id", Value: transfer.AccountDestinationID},
		primitive.E{Key: "amount", Value: transfer.Amount},
		primitive.E{Key: "created_at", Value: transfer.CreatedAt},
	})
	if err != nil {
		return fmt.Errorf("could not create transfer: %w", err)
	}

	return nil
}
