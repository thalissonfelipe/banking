package account

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

var collection *mongo.Collection

func TestMain(m *testing.M) {
	mgoDocker := dockertest.SetupMongoDB()

	collection = mgoDocker.Client.Database("banking").Collection("accounts")

	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"cpf": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("could not create cpf index: %v", err)
	}

	code := m.Run()

	dockertest.RemoveMongoDBContainer(mgoDocker)

	os.Exit(code)
}

func createAccount(t *testing.T, balance int) entities.Account {
	t.Helper()

	acc := entities.NewAccount("Felipe", testdata.GetValidCPF(), testdata.GetValidSecret())
	acc.Balance = balance
	acc.CreatedAt = time.Now()

	_, err := collection.InsertOne(context.Background(), bson.D{
		primitive.E{Key: "id", Value: acc.ID},
		primitive.E{Key: "name", Value: acc.Name},
		primitive.E{Key: "cpf", Value: acc.CPF.String()},
		primitive.E{Key: "secret", Value: acc.Secret.String()},
		primitive.E{Key: "balance", Value: acc.Balance},
		primitive.E{Key: "created_at", Value: primitive.Timestamp{T: uint32(acc.CreatedAt.Unix()), I: 0}},
	})
	require.NoError(t, err)

	return acc
}

func dropCollection(t *testing.T, collection *mongo.Collection) {
	err := collection.Drop(context.Background())
	require.NoError(t, err)
}
