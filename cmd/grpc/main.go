package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/thalissonfelipe/banking/banking/config"
	"github.com/thalissonfelipe/banking/banking/gateway/db/postgres"
	grpcServer "github.com/thalissonfelipe/banking/banking/gateway/grpc"
	"github.com/thalissonfelipe/banking/banking/instrumentation/log"
)

func main() {
	logger := log.New(os.Stderr)

	mainLogger := logger.With(zap.String("module", "main"))

	mainLogger.Info("starting grpc banking application...")

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

	lis, err := net.Listen("tcp", cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("dial connection: %w", err)
	}

	defer lis.Close()

	server := grpcServer.NewServer(logger, conn)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		if serveErr := server.Serve(lis); serveErr != nil {
			mainLogger.Error("failed to serve", zap.Error(serveErr))
		}
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
	<-shutdownCh

	mainLogger.Info("shutting down the server...")

	server.GracefulStop()

	wg.Wait()

	return nil
}
