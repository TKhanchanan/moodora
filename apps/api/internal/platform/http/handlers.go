package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type healthResponse struct {
	Status   string            `json:"status"`
	Checks   map[string]string `json:"checks"`
	Timezone string            `json:"timezone"`
}

type versionResponse struct {
	Name     string `json:"name"`
	Env      string `json:"env"`
	Version  string `json:"version"`
	Timezone string `json:"timezone"`
}

func healthHandler(deps Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		statusCode := http.StatusOK
		status := "ok"
		checks := map[string]string{
			"postgres": "ok",
			"redis":    "ok",
			"storage":  "ok",
		}

		if err := deps.DB.Ping(ctx); err != nil {
			statusCode = http.StatusServiceUnavailable
			status = "degraded"
			checks["postgres"] = "unavailable"
		}

		if err := deps.Redis.Ping(ctx).Err(); err != nil {
			statusCode = http.StatusServiceUnavailable
			status = "degraded"
			checks["redis"] = "unavailable"
		}

		if _, err := deps.Storage.ListBuckets(ctx); err != nil {
			statusCode = http.StatusServiceUnavailable
			status = "degraded"
			checks["storage"] = "unavailable"
		}

		writeJSON(w, statusCode, healthResponse{
			Status:   status,
			Checks:   checks,
			Timezone: deps.Config.TimezoneName(),
		})
	}
}

func versionHandler(deps Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, versionResponse{
			Name:     deps.Config.AppName,
			Env:      deps.Config.AppEnv,
			Version:  "0.1.0",
			Timezone: deps.Config.TimezoneName(),
		})
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
