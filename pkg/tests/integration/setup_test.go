package integration

import (
	"net/http/httptest"
	"os"
	"testing"

	h "github.com/thalissonfelipe/banking/pkg/gateways/http"
	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
)

var (
	pgDocker *dockertest.PostgresDocker
	server   *httptest.Server
)

func TestMain(m *testing.M) {
	pgDocker = dockertest.SetupTest("../../gateways/db/postgres/migrations")
	r := h.NewRouter(pgDocker.DB)
	server = httptest.NewServer(r)

	exitCode := m.Run()

	dockertest.RemoveContainer(pgDocker)
	server.Close()

	os.Exit(exitCode)
}
