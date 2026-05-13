package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/moodora/moodora/apps/api/internal/platform/config"
	"github.com/redis/go-redis/v9"
)

type Dependencies struct {
	Config  config.Config
	DB      *pgxpool.Pool
	Redis   *redis.Client
	Storage *minio.Client
	Logger  *slog.Logger
}

func NewServer(deps Dependencies) *http.Server {
	return &http.Server{
		Addr:              deps.Config.HTTPAddr(),
		Handler:           NewRouter(deps),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
