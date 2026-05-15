package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/moodora/moodora/apps/api/internal/modules/tarot"
	"github.com/moodora/moodora/apps/api/internal/platform/database"
)

func main() {
	sourcePath := flag.String("source", "", "optional tarotapi.dev-compatible JSON source file")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logger.Error("DATABASE_URL is required")
		os.Exit(1)
	}

	ctx := context.Background()
	db, err := database.Connect(ctx, databaseURL)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	var provider tarot.CardProvider = tarot.BuiltInProvider{}
	if *sourcePath != "" {
		provider = tarot.FileProvider{Path: *sourcePath}
	}

	service := tarot.NewService(tarot.NewRepository(db))
	count, err := service.ImportCards(ctx, provider)
	if err != nil {
		logger.Error("failed to import tarot cards", "error", err)
		os.Exit(1)
	}

	logger.Info("tarot cards imported", "count", count)
}
