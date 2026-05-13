package http

import (
	"fmt"
	"net/http"

	"github.com/moodora/moodora/apps/api/internal/modules/checkin"
	"github.com/moodora/moodora/apps/api/internal/modules/lifestyle"
	"github.com/moodora/moodora/apps/api/internal/modules/tarot"
	"github.com/moodora/moodora/apps/api/internal/modules/wallet"
)

func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler(deps))
	mux.HandleFunc("GET /api/v1/version", versionHandler(deps))

	tarotRepo := tarot.NewRepository(deps.DB)
	tarotHandler := tarot.NewHandler(tarot.NewService(tarotRepo))
	tarotHandler.Register(mux)

	walletService := wallet.NewService(deps.DB)
	resolveDevUser := func(r *http.Request) (string, error) {
		if deps.Config.DevUserID == "" {
			return "", fmt.Errorf("DEV_USER_ID must be configured for development wallet endpoints")
		}
		return deps.Config.DevUserID, nil
	}
	wallet.NewHandler(walletService, resolveDevUser).Register(mux)
	checkin.NewHandler(checkin.NewService(deps.DB, walletService, deps.Config.AppTimezone), resolveDevUser).Register(mux)

	lifestyleService := lifestyle.NewService(lifestyle.NewRepository(deps.DB), deps.Config.AppTimezone)
	lifestyle.NewHandler(lifestyleService, func(r *http.Request) string {
		return deps.Config.DevUserID
	}).Register(mux)

	return requestLogger(deps.Logger, mux)
}
