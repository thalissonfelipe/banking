package postgres

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrations embed.FS

func GetMigrationHandler(dbURL string) (*migrate.Migrate, error) {
	const path = "migrations"

	driver, err := iofs.New(migrations, path)
	if err != nil {
		return nil, fmt.Errorf("creating driver from io/fs: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", driver, dbURL)
	if err != nil {
		return nil, fmt.Errorf("creating migrate instance: %w", err)
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
		return fmt.Errorf("applying up migrations: %w", err)
	}

	return nil
}
