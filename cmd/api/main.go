package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/thalissonfelipe/banking/banking/config"
	"github.com/thalissonfelipe/banking/banking/gateway/db/postgres"
	handler "github.com/thalissonfelipe/banking/banking/gateway/http"
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

	mainLogger.Info("starting banking application...")

	cfg, err := config.LoadConfig()
	if err != nil {
		mainLogger.Panic("failed to load config", zap.Error(err))
	}

	if err := startApp(cfg, logger, mainLogger); err != nil {
		mainLogger.Panic("failed to start application", zap.Error(err))
	}
}

func startApp(cfg *config.Config, logger, mainLogger *zap.Logger) error {
	conn, err := pgx.Connect(context.Background(), cfg.Postgres.DSN())
	if err != nil {
		return fmt.Errorf("connecting to postgres: %w", err)
	}

	defer conn.Close(context.Background())

	err = postgres.RunMigrations(cfg.Postgres.DSN())
	if err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	router := handler.NewRouter(logger, conn)
	server := http.Server{
		Handler: router,
		Addr:    cfg.API.Address(),
	}

	go func() {
		if listenErr := server.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			mainLogger.Error("failed to listen and serve", zap.Error(listenErr))
		}
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
	<-shutdownCh

	mainLogger.Info("shuting down the server...")

	const timeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutting down the server: %w", err)
	}

	return nil
}
