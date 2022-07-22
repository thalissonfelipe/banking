package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/thalissonfelipe/banking/banking/config"
	"github.com/thalissonfelipe/banking/banking/gateway/db/postgres"
	handler "github.com/thalissonfelipe/banking/banking/gateway/http"
	"github.com/thalissonfelipe/banking/banking/instrumentation/log"
	"github.com/thalissonfelipe/banking/banking/instrumentation/tracer"
	_ "github.com/thalissonfelipe/banking/docs/swagger"
)

func main() {
	logger := log.New(os.Stderr)

	mainLogger := logger.With(zap.String("module", "main"))

	mainLogger.Info("starting http banking application...")

	cfg, err := config.LoadConfig()
	if err != nil {
		mainLogger.Panic("failed to load config", zap.Error(err))
	}

	if err := startApp(cfg, logger, mainLogger); err != nil {
		mainLogger.Panic("failed to start application", zap.Error(err))
	}
}

func startApp(cfg *config.Config, logger, mainLogger *zap.Logger) error {
	conn, err := pgxpool.Connect(context.Background(), cfg.Postgres.DSN())
	if err != nil {
		return fmt.Errorf("connecting to postgres: %w", err)
	}

	defer conn.Close()

	err = postgres.RunMigrations(cfg.Postgres.DSN())
	if err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	closer, err := tracer.New()
	if err != nil {
		return fmt.Errorf("creating otel tracer: %w", err)
	}

	defer func() {
		if closerErr := closer(); closerErr != nil {
			mainLogger.Error("closing exporter and tracer: %w", zap.Error(closerErr))
		}
	}()

	const (
		readTimeout  = 5 * time.Second
		writeTimeout = 15 * time.Second
	)

	router := handler.NewRouter(logger, conn)
	server := http.Server{
		Handler:           router,
		Addr:              cfg.API.Address(),
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readTimeout,
		WriteTimeout:      writeTimeout,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

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

	err = server.Shutdown(ctx)

	wg.Wait()

	if err != nil {
		return fmt.Errorf("shutting down the server: %w", err)
	}

	return nil
}
