package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"

	h "github.com/thalissonfelipe/banking/pkg/gateways/http"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %s", err.Error())
	}

	conn, err := pgx.ConnectConfig(context.Background(), cfg.Postgres)
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err.Error())
	}
	defer conn.Close(context.Background())

	router := h.NewRouter(conn)

	addr := fmt.Sprintf("%s:%s", cfg.API.Host, cfg.API.Port)
	server := http.Server{
		Handler: router,
		Addr:    addr,
	}

	log.Printf("Server listening on %s!\n", addr)
	log.Fatal(server.ListenAndServe())
}
