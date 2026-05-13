package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moodora/moodora/apps/api/internal/platform/config"
)

func TestCORSMiddlewareAllowsConfiguredOrigin(t *testing.T) {
	handler := corsMiddleware(config.Config{
		AppEnv: "local",
		CORS:   config.CORSConfig{AllowedOrigins: []string{"http://localhost:3000"}},
	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("Access-Control-Allow-Origin = %q", got)
	}
}

func TestCORSMiddlewarePreflightReturnsNoContent(t *testing.T) {
	handler := corsMiddleware(config.Config{
		AppEnv: "local",
		CORS:   config.CORSConfig{AllowedOrigins: []string{"http://localhost:3000"}},
	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not run for OPTIONS preflight")
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/version", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Fatal("Access-Control-Allow-Methods should be set")
	}
}

func TestCORSMiddlewareDoesNotAllowWildcardInProduction(t *testing.T) {
	handler := corsMiddleware(config.Config{
		AppEnv: "production",
		CORS:   config.CORSConfig{AllowedOrigins: []string{"*"}},
	}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want empty", got)
	}
}
