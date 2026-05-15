package checkin

import (
	"net/http"

	"github.com/moodora/moodora/apps/api/internal/modules/wallet"
	"github.com/moodora/moodora/apps/api/internal/platform/httputil"
)

type Handler struct {
	service     *Service
	resolveUser wallet.UserResolver
}

func NewHandler(service *Service, resolveUser wallet.UserResolver) *Handler {
	return &Handler{service: service, resolveUser: resolveUser}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/check-ins", h.create)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	userID, err := h.resolveUser(r)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "dev_user_missing", err.Error())
		return
	}

	response, err := h.service.CheckIn(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "internal_error", "failed to check in")
		return
	}
	httputil.JSON(w, http.StatusOK, response)
}
