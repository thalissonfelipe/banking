package transfer

import (
	"os"
	"testing"

	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
)

var pgDocker *dockertest.PostgresDocker

func TestMain(m *testing.M) {
	pgDocker = dockertest.SetupTest("../migrations")

	exitCode := m.Run()

	dockertest.RemoveContainer(pgDocker)

	os.Exit(exitCode)
}
