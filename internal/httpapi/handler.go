package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"fizz-buzz/internal/service"
	"fizz-buzz/internal/stats"
)

const (
	errorCodeInvalidParameter = "INVALID_PARAMETER"
	errorCodeInternal         = "INTERNAL_ERROR"
)

// Handler groups HTTP dependencies.
type Handler struct {
	stats    *stats.Store
	maxLimit int
}

type errorBody struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewHandler wires the HTTP handlers.
func NewHandler(statsStore *stats.Store, maxLimit int) *Handler {
	return &Handler{stats: statsStore, maxLimit: maxLimit}
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
	params, err := h.parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, errorCodeInvalidParameter, err.Error())
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

func (h *Handler) parseParams(r *http.Request) (service.FizzBuzzParams, error) {
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

	if err := params.Validate(h.maxLimit); err != nil {
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
		http.Error(w, `{"error":{"code":"`+errorCodeInternal+`","message":"failed to encode response"}}`, http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, errorBody{
		Error: apiError{
			Code:    code,
			Message: message,
		},
	})
}
