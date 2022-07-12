package postgres

import (
	"embed"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // driver
	"github.com/golang-migrate/migrate/v4/source/httpfs"
)

//go:embed migrations
var migrations embed.FS

func GetMigrationHandler(dbURL string) (*migrate.Migrate, error) {
	// use httpFS until go-migrate implements ioFS
	// (see https://github.com/golang-migrate/migrate/issues/480#issuecomment-731518493)
	source, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return nil, fmt.Errorf("could not create a migrate source driver from httpfs: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("httpfs", source, dbURL)
	if err != nil {
		return nil, fmt.Errorf("could not create a migrate instance: %w", err)
	}

	return m, nil
}

func RunMigrations(dbURL string) error {
	m, err := GetMigrationHandler(dbURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not apply up migrations: %w", err)
	}

	return nil
}
