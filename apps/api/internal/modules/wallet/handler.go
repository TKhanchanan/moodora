package wallet

import (
	"net/http"

	"github.com/moodora/moodora/apps/api/internal/platform/httputil"
)

type UserResolver func(*http.Request) (string, error)

type Handler struct {
	service     *Service
	resolveUser UserResolver
}

func NewHandler(service *Service, resolveUser UserResolver) *Handler {
	return &Handler{service: service, resolveUser: resolveUser}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/wallet", h.getWallet)
	mux.HandleFunc("GET /api/v1/coin-transactions", h.listTransactions)
}

func (h *Handler) getWallet(w http.ResponseWriter, r *http.Request) {
	userID, err := h.resolveUser(r)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "dev_user_missing", err.Error())
		return
	}
	wallet, err := h.service.GetWallet(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "internal_error", "failed to load wallet")
		return
	}
	httputil.JSON(w, http.StatusOK, wallet)
}

func (h *Handler) listTransactions(w http.ResponseWriter, r *http.Request) {
	userID, err := h.resolveUser(r)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "dev_user_missing", err.Error())
		return
	}
	transactions, err := h.service.ListTransactions(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "internal_error", "failed to list coin transactions")
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"transactions": transactions})
}
