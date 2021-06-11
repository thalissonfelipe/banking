package dockertest

import (
	"context"
	"fmt"
	"log"

	"github.com/ory/dockertest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
