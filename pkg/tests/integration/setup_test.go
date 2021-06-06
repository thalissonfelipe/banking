package integration

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/gateways/hash"
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

func createAccount(t *testing.T, cpf vos.CPF, secret vos.Secret) entities.Account {
	acc := entities.NewAccount("Felipe", cpf, secret)

	encrypter := hash.Hash{}

	err := acc.Secret.Hash(encrypter)
	require.NoError(t, err)

	_, err = pgDocker.DB.Exec(context.Background(),
		`insert into accounts (id, name, cpf, secret, balance) values ($1, $2, $3, $4, $5)`,
		acc.ID, acc.Name, acc.CPF, acc.Secret, acc.Balance)
	require.NoError(t, err)

	return acc
}
