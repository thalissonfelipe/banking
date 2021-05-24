package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"github.com/thalissonfelipe/banking/pkg/config"
	"github.com/thalissonfelipe/banking/pkg/gateways/db/postgres"
	h "github.com/thalissonfelipe/banking/pkg/gateways/http"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("unable to load config")
	}

	conn, err := pgx.Connect(context.Background(), cfg.Postgres.DSN())
	if err != nil {
		log.WithError(err).Fatal("unable to connect to database")
	}
	defer conn.Close(context.Background())

	err = postgres.RunMigrations(cfg.Postgres.DSN())
	if err != nil {
		log.WithError(err).Fatal("unable to run migrations")
	}

	router := h.NewRouter(conn)

	addr := fmt.Sprintf("%s:%s", cfg.API.Host, cfg.API.Port)
	server := http.Server{
		Handler: router,
		Addr:    addr,
	}

	log.Infof("Server listening on %s!", addr)
	log.Fatal(server.ListenAndServe())
}
