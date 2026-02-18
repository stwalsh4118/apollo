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

	"github.com/sean/apollo/api/internal/config"
	"github.com/sean/apollo/api/internal/database"
	"github.com/sean/apollo/api/internal/logging"
	"github.com/sean/apollo/api/internal/server"
)

const (
	exitCodeFailure     = 1
	shutdownGracePeriod = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "apollo startup failed: %v\n", err)
		os.Exit(exitCodeFailure)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger, err := logging.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("create logger: %w", err)
	}

	handle, err := database.Open(ctx, cfg.DatabasePath, logger)
	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}
	defer func() {
		_ = handle.Close()
	}()

	srv := server.New(handle, logger)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:           srv.Router(),
		ReadHeaderTimeout: shutdownGracePeriod,
		IdleTimeout:       2 * time.Minute,
	}

	errCh := make(chan error, 1)

	go func() {
		logger.Info().
			Int("port", cfg.ServerPort).
			Msg("apollo HTTP server starting")

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http server: %w", err)
		}

		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		logger.Info().Msg("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGracePeriod)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	logger.Info().Msg("apollo HTTP server stopped gracefully")

	return nil
}
