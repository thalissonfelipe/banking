package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/thalissonfelipe/banking/banking/config"
	"github.com/thalissonfelipe/banking/banking/gateway/db/postgres"
	h "github.com/thalissonfelipe/banking/banking/gateway/http"
	"github.com/thalissonfelipe/banking/banking/instrumentation/log"
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
	logger := log.New(os.Stderr)

	mainLogger := logger.With(zap.String("module", "main"))

	cfg, err := config.LoadConfig()
	if err != nil {
		mainLogger.Panic("failed to load config", zap.Error(err))
	}

	conn, err := pgx.Connect(context.Background(), cfg.Postgres.DSN())
	if err != nil {
		mainLogger.Panic("failed to connect to database", zap.Error(err))
	}
	defer conn.Close(context.Background())

	err = postgres.RunMigrations(cfg.Postgres.DSN())
	if err != nil {
		mainLogger.Panic("failed to run migrations", zap.Error(err))
	}

	router := h.NewRouter(logger, conn)

	addr := fmt.Sprintf("%s:%s", cfg.API.Host, cfg.API.Port)
	server := http.Server{
		Handler: router,
		Addr:    addr,
	}

	logger.Info("server listening", zap.String("address", addr))

	server.ListenAndServe()
}
