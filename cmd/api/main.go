package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"

	"github.com/thalissonfelipe/banking/banking/config"
	"github.com/thalissonfelipe/banking/banking/gateways/db/postgres"
	h "github.com/thalissonfelipe/banking/banking/gateways/http"
	_ "github.com/thalissonfelipe/banking/docs/swagger"
)

// @title Swagger Banking API
// @version 1.0
// @description This is a simple banking api.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5000
// @BasePath /api/v1
// @query.collection.format multi

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
	if err != nil {
		log.Panicf("unable to run migrations: %v", err)
	}

	router := h.NewRouter(conn)

	addr := fmt.Sprintf("%s:%s", cfg.API.Host, cfg.API.Port)
	server := http.Server{
		Handler: router,
		Addr:    addr,
	}

	log.Printf("Server listening on %s!\n", addr)
	log.Panic(server.ListenAndServe())
}
