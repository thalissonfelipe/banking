package account

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

var collection *mongo.Collection

func TestMain(m *testing.M) {
	mgoDocker := dockertest.SetupMongoDB()

	collection = mgoDocker.Client.Database("banking").Collection("accounts")

	code := m.Run()

	dockertest.RemoveMongoDBContainer(mgoDocker)

	os.Exit(code)
}

func createAccount(t *testing.T, balance int) entities.Account {
	t.Helper()

	acc := entities.NewAccount("Felipe", testdata.GetValidCPF(), testdata.GetValidSecret())
	acc.Balance = balance

	_, err := collection.InsertOne(context.Background(), bson.D{
		primitive.E{Key: "id", Value: acc.ID},
		primitive.E{Key: "name", Value: acc.Name},
		primitive.E{Key: "cpf", Value: acc.CPF.String()},
		primitive.E{Key: "secret", Value: acc.Secret.String()},
		primitive.E{Key: "balance", Value: acc.Balance},
	})
	require.NoError(t, err)

	return acc
}
