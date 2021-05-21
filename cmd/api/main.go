package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"

	"github.com/thalissonfelipe/banking/pkg/config"
	"github.com/thalissonfelipe/banking/pkg/gateways/db/postgres"
	h "github.com/thalissonfelipe/banking/pkg/gateways/http"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %s", err.Error())
	}

	conn, err := pgx.Connect(context.Background(), cfg.Postgres.DSN())
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err.Error())
	}
	defer conn.Close(context.Background())

	err = postgres.RunMigrations(cfg.Postgres.DSN())
	log.Println(cfg.Postgres.DSN())
	if err != nil {
		log.Fatalf("unable to run migrations: %s", err.Error())
	}

	router := h.NewRouter(conn)

	addr := fmt.Sprintf("%s:%s", cfg.API.Host, cfg.API.Port)
	server := http.Server{
		Handler: router,
		Addr:    addr,
	}

	log.Printf("Server listening on %s!\n", addr)
	log.Fatal(server.ListenAndServe())
}
