package transfer

import (
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
)

var db *mongo.Database

func TestMain(m *testing.M) {
	mgoDocker := dockertest.SetupMongoDB()

	db = mgoDocker.Client.Database("banking")

	code := m.Run()

	dockertest.RemoveMongoDBContainer(mgoDocker)

	os.Exit(code)
}
