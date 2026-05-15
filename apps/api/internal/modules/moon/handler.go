package moon

import (
	"errors"
	"net/http"
	"strings"

	"github.com/moodora/moodora/apps/api/internal/platform/httputil"
)

type UserResolver func(*http.Request) string

type Handler struct {
	service     *Service
	resolveUser UserResolver
}

func NewHandler(service *Service, resolveUser UserResolver) *Handler {
	return &Handler{service: service, resolveUser: resolveUser}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/moon/today", h.today)
	mux.HandleFunc("POST /api/v1/moon/birthday", h.birthday)
	mux.HandleFunc("GET /api/v1/moon/reports/{id}", h.getReport)
}

func (h *Handler) today(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.Today(r.Context(), h.resolveUser(r))
	writeMoon(w, report, err)
}

func (h *Handler) birthday(w http.ResponseWriter, r *http.Request) {
	var req BirthdayRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
		return
	}
	report, err := h.service.Birthday(r.Context(), h.resolveUser(r), req)
	writeMoon(w, report, err)
}

func (h *Handler) getReport(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	report, err := h.service.GetReport(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "moon report not found")
		return
	}
	writeMoon(w, report, err)
}

func writeMoon(w http.ResponseWriter, report ReportResponse, err error) {
	if err != nil {
		switch err.Error() {
		case "birthDate is required", "date must use YYYY-MM-DD", "invalid timezone":
			httputil.Error(w, http.StatusBadRequest, "invalid_request", err.Error())
		default:
			httputil.Error(w, http.StatusInternalServerError, "internal_error", "moon request failed")
		}
		return
	}
	httputil.JSON(w, http.StatusOK, report)
}
