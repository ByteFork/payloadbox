package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGzip_EncodesWhenAccepted(t *testing.T) {
	payload := strings.Repeat("hello-gzip-world ", 200)
	h := Gzip(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = io.WriteString(w, payload)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", http.NoBody)
	req.Header.Set("Accept-Encoding", "gzip")
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("Content-Encoding = %q, want gzip", got)
	}

	gr, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatalf("body is not valid gzip: %v", err)
	}

	out, err := io.ReadAll(gr)
	if err != nil {
		t.Fatal(err)
	}

	if string(out) != payload {
		t.Errorf("decompressed len=%d, want len=%d", len(out), len(payload))
	}
}

func TestGzip_PassesThroughWhenNotAccepted(t *testing.T) {
	payload := "plain response"
	h := Gzip(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, payload)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", http.NoBody)
	// No Accept-Encoding header.
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Content-Encoding"); got != "" {
		t.Errorf("Content-Encoding = %q, want empty", got)
	}

	if rec.Body.String() != payload {
		t.Errorf("body = %q, want %q", rec.Body.String(), payload)
	}
}

func TestGzip_PassesThroughWhenOtherEncoding(t *testing.T) {
	h := Gzip(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "x")
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", http.NoBody)
	req.Header.Set("Accept-Encoding", "br, deflate") // no gzip
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Content-Encoding"); got != "" {
		t.Errorf("Content-Encoding = %q, want empty", got)
	}
}
