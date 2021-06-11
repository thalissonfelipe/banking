package account

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
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

func createAccount(t *testing.T, r Repository, balance int) entities.Account {
	t.Helper()

	acc := entities.NewAccount("Felipe", testdata.GetValidCPF(), testdata.GetValidSecret())
	acc.Balance = balance

	err := r.CreateAccount(context.Background(), &acc)
	require.NoError(t, err)

	return acc
}
