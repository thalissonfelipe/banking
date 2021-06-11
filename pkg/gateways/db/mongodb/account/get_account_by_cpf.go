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

func (r Repository) GetAccountByCPF(ctx context.Context, cpf vos.CPF) (*entities.Account, error) {
	var account entities.Account

	err := r.collection.FindOne(ctx, bson.M{"cpf": cpf.String()}).Decode(&account)
	if err == nil {
		return &account, nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, entities.ErrAccountDoesNotExist
	}

	return nil, fmt.Errorf("could not get account by cpf: %w", err)
}
