package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS_NoOriginPassesThrough(t *testing.T) {
	called := false
	h := WithCORS(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", http.NoBody)
	// No Origin header.
	h.ServeHTTP(rec, req)

	if !called {
		t.Fatal("next handler was not called")
	}

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Allow-Origin = %q without Origin header, want empty", got)
	}
}

func TestCORS_SimpleRequestAddsHeaders(t *testing.T) {
	called := false
	h := WithCORS(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", http.NoBody)
	req.Header.Set("Origin", "https://example.com")
	h.ServeHTTP(rec, req)

	if !called {
		t.Fatal("next handler was not called")
	}

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Allow-Origin = %q, want *", got)
	}

	if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Error("Allow-Methods missing")
	}

	if got := rec.Header().Get("Vary"); got == "" {
		t.Error("Vary missing")
	}
}

func TestCORS_PreflightShortCircuits(t *testing.T) {
	called := false
	h := WithCORS(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		called = true
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodOptions, "/", http.NoBody)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	h.ServeHTTP(rec, req)

	if called {
		t.Error("next handler was called on preflight; should short-circuit")
	}

	if rec.Code != http.StatusNoContent {
		t.Errorf("preflight status = %d, want 204", rec.Code)
	}

	if rec.Header().Get("Access-Control-Max-Age") == "" {
		t.Error("Access-Control-Max-Age missing on preflight")
	}
}

func TestCORS_OptionsWithoutRequestMethodPassesThrough(t *testing.T) {
	called := false
	h := WithCORS(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		called = true
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodOptions, "/", http.NoBody)
	req.Header.Set("Origin", "https://example.com")
	// Intentionally no Access-Control-Request-Method.
	h.ServeHTTP(rec, req)

	if !called {
		t.Fatal("next handler not called for non-preflight OPTIONS")
	}
}
