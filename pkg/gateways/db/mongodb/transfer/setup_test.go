package transfer

import (
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
)

var collection *mongo.Collection

func TestMain(m *testing.M) {
	mgoDocker := dockertest.SetupMongoDB()

	collection = mgoDocker.Client.Database("banking").Collection("transfers")

	code := m.Run()

	dockertest.RemoveMongoDBContainer(mgoDocker)

	os.Exit(code)
}
