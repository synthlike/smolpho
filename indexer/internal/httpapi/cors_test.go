package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSAllowedOrigin(t *testing.T) {
	handler := WithCORS(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), []string{"http://localhost:5173"})

	request := httptest.NewRequest(http.MethodGet, "/api/v1/config", nil)
	request.Header.Set("Origin", "http://localhost:5173")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("Access-Control-Allow-Origin = %q", got)
	}
}

func TestCORSDoesNotAllowUnconfiguredOrigin(t *testing.T) {
	handler := WithCORS(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), []string{"http://localhost:5173"})

	request := httptest.NewRequest(http.MethodGet, "/api/v1/config", nil)
	request.Header.Set("Origin", "https://example.com")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want empty", got)
	}
}

func TestCORSPreflight(t *testing.T) {
	handler := WithCORS(http.NotFoundHandler(), []string{"http://localhost:5173"})
	request := httptest.NewRequest(http.MethodOptions, "/api/v1/config", nil)
	request.Header.Set("Origin", "http://localhost:5173")
	request.Header.Set("Access-Control-Request-Method", http.MethodGet)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", response.Code)
	}
	if got := response.Header().Get("Access-Control-Allow-Methods"); got != http.MethodGet {
		t.Fatalf("Access-Control-Allow-Methods = %q", got)
	}
}

func TestValidateAllowedOrigins(t *testing.T) {
	for _, origin := range []string{
		"localhost:5173",
		"http://localhost:5173/",
		"https://example.com/app",
		"ftp://example.com",
	} {
		if err := ValidateAllowedOrigins([]string{origin}); err == nil {
			t.Errorf("ValidateAllowedOrigins(%q) succeeded", origin)
		}
	}
	if err := ValidateAllowedOrigins([]string{"http://localhost:5173", "https://app.example.com"}); err != nil {
		t.Fatalf("valid origins rejected: %v", err)
	}
}
