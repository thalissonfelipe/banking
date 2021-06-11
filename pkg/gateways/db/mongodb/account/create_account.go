package account

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (r Repository) CreateAccount(ctx context.Context, account *entities.Account) error {
	createdAt := time.Now()

	_, err := r.db.Collection("accounts").InsertOne(ctx, bson.D{
		primitive.E{Key: "id", Value: account.ID},
		primitive.E{Key: "name", Value: account.Name},
		primitive.E{Key: "cpf", Value: account.CPF.String()},
		primitive.E{Key: "secret", Value: account.Secret.String()},
		primitive.E{Key: "balance", Value: account.Balance},
		primitive.E{Key: "created_at", Value: createdAt},
	})
	if err == nil {
		account.CreatedAt = createdAt

		return nil
	}

	if mongo.IsDuplicateKeyError(err) {
		return entities.ErrAccountAlreadyExists
	}

	return fmt.Errorf("could not create account: %w", err)
}
