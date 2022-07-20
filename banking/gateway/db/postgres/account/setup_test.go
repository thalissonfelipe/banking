package account

import (
	"log"
	"os"
	"testing"

	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	teardownFn, err := dockertest.NewPostgresContainer()
	if err != nil {
		log.Panicf("failed to setup postgres container: %v", err)
	}

	defer teardownFn()

	return m.Run()
}
