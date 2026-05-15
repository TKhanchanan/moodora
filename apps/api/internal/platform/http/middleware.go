package http

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/moodora/moodora/apps/api/internal/platform/config"
)

const allowedCORSMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
const allowedCORSHeaders = "Content-Type, Authorization, X-API-Key"

func corsMiddleware(cfg config.Config, next http.Handler) http.Handler {
	allowedOrigins := map[string]bool{}
	for _, origin := range cfg.CORS.AllowedOrigins {
		allowedOrigins[origin] = true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && isOriginAllowed(origin, allowedOrigins, cfg.AppEnv) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", allowedCORSMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowedCORSHeaders)
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isOriginAllowed(origin string, allowedOrigins map[string]bool, appEnv string) bool {
	if allowedOrigins[origin] {
		return true
	}
	if allowedOrigins["*"] && strings.EqualFold(appEnv, "production") {
		return false
	}
	return allowedOrigins["*"]
}

func requestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		logger.Info(
			"http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.statusCode,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
