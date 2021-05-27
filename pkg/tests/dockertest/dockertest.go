package dockertest

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/ory/dockertest"
)

type PostgresDocker struct {
	DB       *pgx.Conn
	Pool     *dockertest.Pool
	Resource *dockertest.Resource
}

func SetupTest(migrationsPath string) *PostgresDocker {
	var conn *pgx.Conn

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err.Error())
	}

	database := getRandomDBName()

	resource, err := pool.Run(
		"postgres",
		"13.2",
		[]string{"POSTGRES_PASSWORD=postgres", "POSTGRES_DB=" + database},
	)
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	connString := fmt.Sprintf(
		"postgres://postgres:postgres@localhost:%s/%s?sslmode=disable",
		resource.GetPort("5432/tcp"),
		database)

	if err = pool.Retry(func() error {
		var err error
		ctx := context.Background()
		conn, err = pgx.Connect(ctx, connString)
		if err != nil {
			return err
		}
		return conn.Ping(ctx)
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	if err := runMigrations(migrationsPath, connString); err != nil {
		log.Fatalf("could not run migrations: %s", err.Error())
	}

	return &PostgresDocker{
		DB:       conn,
		Pool:     pool,
		Resource: resource,
	}
}

func RemoveContainer(pgDocker *PostgresDocker) {
	if err := pgDocker.Pool.Purge(pgDocker.Resource); err != nil {
		log.Fatalf("could not purge resource: %s", err.Error())
	}
}

func TruncateTables(ctx context.Context, db *pgx.Conn) {
	if _, err := db.Exec(ctx, "truncate transfers, accounts"); err != nil {
		log.Fatalf("could not truncate tables: %s", err.Error())
	}
}

func runMigrations(migrationsPath, connString string) error {
	if migrationsPath != "" {
		mig, err := migrate.New("file://"+migrationsPath, connString)
		if err != nil {
			return fmt.Errorf("failed to start migrate struct: %s", err.Error())
		}
		defer mig.Close()
		if err = mig.Up(); err != nil {
			return fmt.Errorf("failed to run migration: %s", err.Error())
		}
	}

	return nil
}

func getRandomDBName() string {
	return fmt.Sprintf("db%d", rand.NewSource(int64(time.Now().Nanosecond())).Int63())
}
