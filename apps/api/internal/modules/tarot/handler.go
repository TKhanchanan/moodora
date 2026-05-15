package tarot

import (
	"errors"
	"net/http"
	"strings"

	"github.com/moodora/moodora/apps/api/internal/platform/httputil"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/tarot/cards", h.listCards)
	mux.HandleFunc("GET /api/v1/tarot/cards/{sourceCode}", h.getCard)
	mux.HandleFunc("GET /api/v1/tarot/spreads", h.listSpreads)
	mux.HandleFunc("GET /api/v1/tarot/spreads/{code}", h.getSpread)
	mux.HandleFunc("POST /api/v1/tarot/readings", h.createReading)
	mux.HandleFunc("GET /api/v1/tarot/readings/{id}", h.getReading)
}

func (h *Handler) listCards(w http.ResponseWriter, r *http.Request) {
	cards, err := h.service.ListCards(r.Context())
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "internal_error", "failed to list tarot cards")
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"cards": cards})
}

func (h *Handler) getCard(w http.ResponseWriter, r *http.Request) {
	sourceCode := strings.TrimSpace(r.PathValue("sourceCode"))
	card, err := h.service.GetCard(r.Context(), sourceCode)
	if err != nil {
		writeTarotError(w, err, "tarot card not found")
		return
	}
	httputil.JSON(w, http.StatusOK, card)
}

func (h *Handler) listSpreads(w http.ResponseWriter, r *http.Request) {
	spreads, err := h.service.ListSpreads(r.Context())
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "internal_error", "failed to list tarot spreads")
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"spreads": spreads})
}

func (h *Handler) getSpread(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.PathValue("code"))
	spread, err := h.service.GetSpread(r.Context(), code)
	if err != nil {
		writeTarotError(w, err, "tarot spread not found")
		return
	}
	httputil.JSON(w, http.StatusOK, spread)
}

func (h *Handler) createReading(w http.ResponseWriter, r *http.Request) {
	var req CreateReadingRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
		return
	}

	reading, err := h.service.CreateReading(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "tarot spread not found")
			return
		}
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "required") {
			httputil.Error(w, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		httputil.Error(w, http.StatusInternalServerError, "internal_error", "failed to create tarot reading")
		return
	}

	httputil.JSON(w, http.StatusCreated, reading)
}

func (h *Handler) getReading(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	reading, err := h.service.GetReading(r.Context(), id)
	if err != nil {
		writeTarotError(w, err, "tarot reading not found")
		return
	}
	httputil.JSON(w, http.StatusOK, reading)
}

func writeTarotError(w http.ResponseWriter, err error, notFoundMessage string) {
	if errors.Is(err, ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", notFoundMessage)
		return
	}
	httputil.Error(w, http.StatusInternalServerError, "internal_error", "tarot request failed")
}
