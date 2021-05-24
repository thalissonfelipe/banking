package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/thalissonfelipe/banking/pkg/config"
	"github.com/thalissonfelipe/banking/pkg/gateways/db/postgres"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
}

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
								log.WithError(err).Fatal("unable to load config from cli")
							}
							m, err := postgres.GetMigrationHandler(cfg.Postgres.DSN())
							if err != nil {
								log.WithError(err).Fatal("unable to run migrations from cli")
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
								log.WithError(err).Fatal("unable to load config cli")
							}
							m, err := postgres.GetMigrationHandler(cfg.Postgres.DSN())
							if err != nil {
								log.WithError(err).Fatal("unable to run migrations from cli")
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
		log.WithError(err).Fatal("unable to run app")
	}
}
