package main

import (
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
								log.Fatalf("unable to load config: %s", err.Error())
							}
							m, err := postgres.GetMigrationHandler(cfg.Postgres.DSN())
							if err != nil {
								log.Fatal(err)
							}
							return m.Up()
						},
					},
					{
						Name:  "down",
						Usage: "migrate down",
						Action: func(c *cli.Context) error {
							cfg, err := config.LoadConfig()
							if err != nil {
								log.Fatalf("unable to load config: %s", err.Error())
							}
							m, err := postgres.GetMigrationHandler(cfg.Postgres.DSN())
							if err != nil {
								log.Fatal(err)
							}
							return m.Down()
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
