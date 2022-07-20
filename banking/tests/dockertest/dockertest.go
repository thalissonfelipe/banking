package dockertest

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/gateway/db/postgres"
)

var db *pgxpool.Pool

func NewPostgresContainer() (teardownFn func(), err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("creating pool: %w", err)
	}

	database := getRandomDBName()

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.2",
		Env:        []string{"POSTGRES_PASSWORD=postgres", "POSTGRES_DB=" + database},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
	})
	if err != nil {
		return nil, fmt.Errorf("starting docker container: %w", err)
	}

	connString := fmt.Sprintf(
		"postgres://postgres:postgres@localhost:%s/%s?sslmode=disable",
		resource.GetPort("5432/tcp"),
		database)

	if err = pool.Retry(func() error {
		ctx := context.Background()

		db, err = pgxpool.Connect(ctx, connString)
		if err != nil {
			return fmt.Errorf("connecting to postgres: %w", err)
		}

		err = db.Ping(ctx)
		if err != nil {
			return fmt.Errorf("ping connection: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	if err = postgres.RunMigrations(connString); err != nil {
		return nil, fmt.Errorf("running migrations")
	}

	teardownFn = func() {
		resource.Close()
	}

	return teardownFn, nil
}

func NewDB(t *testing.T, name string) *pgxpool.Pool {
	t.Helper()

	if name == "" {
		require.FailNow(t, "name must no be an empty string")
	}

	re := regexp.MustCompile(`[\W]`)

	name = re.ReplaceAllString(strings.ToLower(name), "_")
	dropDatabaseQuery := fmt.Sprintf("drop database if exists %s", name)

	_, err := db.Exec(context.Background(), dropDatabaseQuery)
	require.NoError(t, err)

	_, err = db.Exec(context.Background(), fmt.Sprintf("create database %s", name))
	require.NoError(t, err)

	connString := strings.Replace(db.Config().ConnString(), db.Config().ConnConfig.Database, name, 1)

	newDB, err := pgxpool.Connect(context.Background(), connString)
	require.NoError(t, err)

	err = postgres.RunMigrations(connString)
	require.NoError(t, err)

	t.Cleanup(func() {
		newDB.Close()
		_, err := db.Exec(context.Background(), dropDatabaseQuery)
		require.NoError(t, err)
	})

	return newDB
}

func getRandomDBName() string {
	return fmt.Sprintf("db%d", rand.NewSource(int64(time.Now().Nanosecond())).Int63())
}
