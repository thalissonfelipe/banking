package account

import (
	"context"
	"log"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
)

var db *mongo.Database

func TestMain(m *testing.M) {
	mgoDocker := dockertest.SetupMongoDB()

	db = mgoDocker.Client.Database("banking")

	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"cpf": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := db.Collection("accounts").Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("could not create cpf index: %v", err)
	}

	code := m.Run()

	dockertest.RemoveMongoDBContainer(mgoDocker)

	os.Exit(code)
}

func assertAccountResponse(t *testing.T, want, got entities.Account, ignoreSecret bool) {
	t.Helper()

	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.CPF, got.CPF)
	assert.Equal(t, want.Balance, got.Balance)
	assert.Equal(t, want.CreatedAt.Unix(), got.CreatedAt.Unix())

	if !ignoreSecret {
		assert.Equal(t, want.Secret, got.Secret)
	}
}
