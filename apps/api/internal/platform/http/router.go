package http

import "net/http"

func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler(deps))
	mux.HandleFunc("GET /api/v1/version", versionHandler(deps))

	return requestLogger(deps.Logger, mux)
}
