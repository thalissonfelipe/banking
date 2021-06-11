package account

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func (r Repository) GetAccountByID(ctx context.Context, id vos.ID) (*entities.Account, error) {
	var account entities.Account

	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&account)
	if err == nil {
		return &account, nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, entities.ErrAccountDoesNotExist
	}

	return nil, fmt.Errorf("could not get account by id: %w", err)
}
