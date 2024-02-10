package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/argatu/boilerplate/internal/config"
	"github.com/argatu/boilerplate/internal/logger"
	"github.com/argatu/boilerplate/internal/postgres"
)

var (
	readTimeout     = 5 * time.Second
	writeTimeout    = 5 * time.Second
	idleTimeout     = 120 * time.Second
	shutdownTimeout = 5 * time.Second
)

func Run() error {
	ctx := context.Background()

	log := logger.NewWithEvents(
		os.Stdout,
		logger.LevelInfo,
		"BOILERPLATE",
		func(ctx context.Context) string {
			return traceID(ctx)
		},
		logger.Events{},
	)

	cfg, err := config.New()
	if err != nil {
		return err
	}

	db, err := postgres.Connect(cfg.DSN())
	if err != nil {
		return err
	}
	defer db.Close()

	h := newHandler(log)

	srv := http.Server{
		Addr:         cfg.Addr(),
		Handler:      h.routes(),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		log.Info(ctx, "server listening", "addr", cfg.Server.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Warn(ctx, "server is shutting down", "signal", sig.String())
		defer log.Warn(ctx, "server stopped")

		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			if err := srv.Close(); err != nil {
				return err
			}

			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
