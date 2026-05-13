package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/moodora/moodora/apps/api/internal/platform/cache"
	"github.com/moodora/moodora/apps/api/internal/platform/config"
	"github.com/moodora/moodora/apps/api/internal/platform/database"
	httpserver "github.com/moodora/moodora/apps/api/internal/platform/http"
	"github.com/moodora/moodora/apps/api/internal/platform/storage"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	redisClient, err := cache.Connect(ctx, cfg.RedisURL)
	if err != nil {
		logger.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	storageClient, err := storage.Connect(cfg.Storage)
	if err != nil {
		logger.Error("failed to connect to object storage", "error", err)
		os.Exit(1)
	}

	deps := httpserver.Dependencies{
		Config:  cfg,
		DB:      db,
		Redis:   redisClient,
		Storage: storageClient,
		Logger:  logger,
	}

	server := httpserver.NewServer(deps)

	errCh := make(chan error, 1)
	go func() {
		logger.Info("api server started", "addr", server.Addr, "env", cfg.AppEnv)
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("api server failed", "error", err)
			os.Exit(1)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("api server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("api server stopped")
}
