package dockertest

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

type MongoDBDocker struct {
	Client   *mongo.Client
	Pool     *dockertest.Pool
	Resource *dockertest.Resource
}

func SetupMongoDB() *MongoDBDocker {
	var client *mongo.Client

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %v", err)
	}

	resource, err := pool.Run("mongo", "latest", []string{})
	if err != nil {
		log.Fatalf("could not start resource: %v", err)
	}

	connString := fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))

	// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet.
	if err = pool.Retry(func() error {
		ctx := context.Background()

		client, err = mongo.Connect(ctx, options.Client().ApplyURI(connString))
		if err != nil {
			return fmt.Errorf("could not connect with mongodb: %w", err)
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			return fmt.Errorf("could not ping: %w", err)
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %v", err)
	}

	return &MongoDBDocker{
		Client:   client,
		Pool:     pool,
		Resource: resource,
	}
}

func RemoveMongoDBContainer(mgoDocker *MongoDBDocker) {
	if err := mgoDocker.Pool.Purge(mgoDocker.Resource); err != nil {
		log.Fatalf("Could not purge resource: %v", err)
	}
}

func DropCollection(t *testing.T, collection *mongo.Collection) {
	err := collection.Drop(context.Background())
	require.NoError(t, err)
}

func CreateAccount(t *testing.T, collection *mongo.Collection, balance int) entities.Account {
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
		primitive.E{Key: "created_at", Value: acc.CreatedAt},
	})
	require.NoError(t, err)

	return acc
}
