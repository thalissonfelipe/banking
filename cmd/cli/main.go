package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/thalissonfelipe/banking/pkg/config"
	"github.com/thalissonfelipe/banking/pkg/gateways/db/postgres"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "migrate",
				Aliases: []string{"m"},
				Usage:   "migrate commands",
				Subcommands: []*cli.Command{
					{
						Name:  "up",
						Usage: "migrate up",
						Action: func(c *cli.Context) error {
							cfg, err := config.LoadConfig()
							if err != nil {
								return fmt.Errorf("could not load config: %w", err)
							}

							m, err := postgres.GetMigrationHandler(cfg.Postgres.DSN())
							if err != nil {
								return fmt.Errorf("could not get migration handler: %w", err)
							}

							err = m.Up()
							if err != nil {
								return fmt.Errorf("could not apply up migrations: %w", err)
							}

							return nil
						},
					},
					{
						Name:  "down",
						Usage: "migrate down",
						Action: func(c *cli.Context) error {
							cfg, err := config.LoadConfig()
							if err != nil {
								return fmt.Errorf("could not load config: %w", err)
							}

							m, err := postgres.GetMigrationHandler(cfg.Postgres.DSN())
							if err != nil {
								return fmt.Errorf("could not get migration handler: %w", err)
							}

							err = m.Down()
							if err != nil {
								return fmt.Errorf("could not apply up migrations: %w", err)
							}

							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("could not run app: %v", err)
	}
}
