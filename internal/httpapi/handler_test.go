package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fizz-buzz/internal/stats"
)

func TestFizzBuzzEndpoint(t *testing.T) {
	handler := NewHandler(stats.NewStore()).Routes()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got []string
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	want := []string{
		"1", "2", "fizz", "4", "buzz",
		"fizz", "7", "8", "fizz", "buzz",
		"11", "fizz", "13", "14", "fizzbuzz",
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d items, got %d", len(want), len(got))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("item %d: expected %q, got %q", i, want[i], got[i])
		}
	}
}

func TestFizzBuzzEndpointValidationError(t *testing.T) {
	handler := NewHandler(stats.NewStore()).Routes()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/fizzbuzz?int1=0&int2=5&limit=15&str1=fizz&str2=buzz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestStatisticsEndpoint(t *testing.T) {
	handler := NewHandler(stats.NewStore()).Routes()

	requests := []string{
		"/api/v1/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz",
		"/api/v1/fizzbuzz?int1=2&int2=7&limit=14&str1=foo&str2=bar",
		"/api/v1/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz",
	}

	for _, path := range requests {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/statistics", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got struct {
		Params struct {
			Int1  int    `json:"int1"`
			Int2  int    `json:"int2"`
			Limit int    `json:"limit"`
			Str1  string `json:"str1"`
			Str2  string `json:"str2"`
		} `json:"params"`
		Hits int `json:"hits"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.Hits != 2 {
		t.Fatalf("expected hits 2, got %d", got.Hits)
	}

	if got.Params.Int1 != 3 || got.Params.Int2 != 5 || got.Params.Limit != 15 || got.Params.Str1 != "fizz" || got.Params.Str2 != "buzz" {
		t.Fatalf("unexpected top params: %+v", got.Params)
	}
}
