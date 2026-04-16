package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"fizz-buzz/internal/service"
	"fizz-buzz/internal/stats"
)

// Handler groups HTTP dependencies.
type Handler struct {
	stats *stats.Store
}

// NewHandler wires the HTTP handlers.
func NewHandler(statsStore *stats.Store) *Handler {
	return &Handler{stats: statsStore}
}

// Routes returns the API router.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/fizzbuzz", h.handleFizzBuzz)
	mux.HandleFunc("GET /api/v1/statistics", h.handleStatistics)
	mux.HandleFunc("GET /health", h.handleHealth)
	return mux
}

func (h *Handler) handleFizzBuzz(w http.ResponseWriter, r *http.Request) {
	params, err := parseParams(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	result := service.FizzBuzz(params)
	// Only successful requests are recorded in the statistics store.
	h.stats.Record(params)

	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) handleStatistics(w http.ResponseWriter, _ *http.Request) {
	entry, ok := h.stats.Top()
	if !ok {
		writeJSON(w, http.StatusOK, map[string]any{
			"params": nil,
			"hits":   0,
		})
		return
	}

	writeJSON(w, http.StatusOK, entry)
}

func (h *Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func parseParams(r *http.Request) (service.FizzBuzzParams, error) {
	query := r.URL.Query()

	// Parse numeric fields first so the API can return precise validation errors.
	int1, err := parsePositiveInt(query.Get("int1"), "int1")
	if err != nil {
		return service.FizzBuzzParams{}, err
	}

	int2, err := parsePositiveInt(query.Get("int2"), "int2")
	if err != nil {
		return service.FizzBuzzParams{}, err
	}

	limit, err := parsePositiveInt(query.Get("limit"), "limit")
	if err != nil {
		return service.FizzBuzzParams{}, err
	}

	params := service.FizzBuzzParams{
		Int1:  int1,
		Int2:  int2,
		Limit: limit,
		Str1:  query.Get("str1"),
		Str2:  query.Get("str2"),
	}

	if err := params.Validate(); err != nil {
		return service.FizzBuzzParams{}, err
	}

	return params, nil
}

func parsePositiveInt(value string, field string) (int, error) {
	if value == "" {
		return 0, errors.New(field + " is required")
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New(field + " must be a valid integer")
	}

	return parsed, nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	// Encoding should normally succeed with our payloads, but we still fail safely.
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}
