package lifestyle

import (
	"net/http"

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
	mux.HandleFunc("GET /api/v1/lucky-colors/today", h.todayColors)
	mux.HandleFunc("GET /api/v1/lucky-foods/today", h.todayFoods)
	mux.HandleFunc("GET /api/v1/lucky-items/today", h.todayItems)
	mux.HandleFunc("GET /api/v1/avoidance/today", h.todayAvoidance)
	mux.HandleFunc("GET /api/v1/daily-insights/today", h.todayDailyInsight)
}

func (h *Handler) todayColors(w http.ResponseWriter, r *http.Request) {
	insight, err := h.service.TodayColors(r.Context(), h.resolveUser(r), r.URL.Query().Get("purpose"))
	writeLifestyle(w, insight, err)
}

func (h *Handler) todayFoods(w http.ResponseWriter, r *http.Request) {
	insight, err := h.service.TodayFoods(r.Context(), h.resolveUser(r))
	writeLifestyle(w, insight, err)
}

func (h *Handler) todayItems(w http.ResponseWriter, r *http.Request) {
	insight, err := h.service.TodayItems(r.Context(), h.resolveUser(r))
	writeLifestyle(w, insight, err)
}

func (h *Handler) todayAvoidance(w http.ResponseWriter, r *http.Request) {
	insight, err := h.service.TodayAvoidances(r.Context(), h.resolveUser(r))
	writeLifestyle(w, insight, err)
}

func (h *Handler) todayDailyInsight(w http.ResponseWriter, r *http.Request) {
	insight, err := h.service.DailyInsight(r.Context(), h.resolveUser(r), r.URL.Query().Get("purpose"))
	writeLifestyle(w, insight, err)
}

func writeLifestyle(w http.ResponseWriter, insight DailyInsight, err error) {
	if err != nil {
		if err.Error() == "invalid purpose" {
			httputil.Error(w, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		httputil.Error(w, http.StatusInternalServerError, "internal_error", "lifestyle request failed")
		return
	}
	httputil.JSON(w, http.StatusOK, insight)
}
