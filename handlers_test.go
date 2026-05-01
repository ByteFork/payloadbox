package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ByteFork/payloadbox/internal/middleware"
)

// newTestHandler builds the mux wired with the server handlers, wrapped with CORS.
// Most tests dispatch requests in-process via handler.ServeHTTP, which is faster
// than a real httptest.Server and still exercises the full middleware + mux chain.
func newTestHandler(t *testing.T, settings ServerSettings) (http.Handler, *Server) {
	t.Helper()

	server := NewServer(settings)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.Health)
	mux.HandleFunc("/version", server.Version)
	mux.HandleFunc("/api/v1/settings", server.Settings)
	mux.Handle("/api/v1/history", middleware.Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			server.ClearRecords(w, r)
			return
		}

		server.ListRecords(w, r)
	})))
	mux.Handle("GET /api/v1/history/{id}", middleware.Gzip(http.HandlerFunc(server.GetRecord)))
	mux.HandleFunc("/api/v1/events", server.Events)
	mux.HandleFunc("/", server.Record)

	return middleware.WithCORS(mux), server
}

// newTestServer wraps newTestHandler in a real httptest.Server.
// Use only for tests that need streaming HTTP, such as SSE, since httptest.NewRecorder
// buffers the entire response and cannot exercise flushes.
func newTestServer(t *testing.T, settings ServerSettings) (*httptest.Server, *Server) {
	t.Helper()

	handler, server := newTestHandler(t, settings)
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	return ts, server
}

func defaultSettings() ServerSettings {
	return ServerSettings{
		Address:           ":0",
		MaxBodySizeBytes:  1024,
		MaxRecordsToStore: 100,
		LogRequests:       false, // silent tests
		LogLevel:          "info",
	}
}

func TestRecord_CapturesAndStores(t *testing.T) {
	handler, s := newTestHandler(t, defaultSettings())

	body := `{"hello":"world"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/webhooks/foo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202", rec.Code)
	}

	records := s.store.List()
	if len(records) != 1 {
		t.Fatalf("store has %d records, want 1", len(records))
	}

	got := records[0]
	if got.ID == "" {
		t.Fatal("record has empty ID")
	}

	if headerID := rec.Header().Get(payloadboxRecordIDHeader); headerID != got.ID {
		t.Errorf("%s = %q, want stored record ID %q", payloadboxRecordIDHeader, headerID, got.ID)
	}

	if got.Request.Method != http.MethodPost {
		t.Errorf("method = %q, want POST", got.Request.Method)
	}

	if got.Request.Path != "/webhooks/foo" {
		t.Errorf("path = %q, want /webhooks/foo", got.Request.Path)
	}

	if got.Request.Body != body {
		t.Errorf("body = %q, want %q", got.Request.Body, body)
	}

	if got.Response.StatusCode != http.StatusAccepted {
		t.Errorf("response status = %d, want 202", got.Response.StatusCode)
	}

	if gotHeader := http.Header(got.Response.Headers).Get(payloadboxRecordIDHeader); gotHeader != got.ID {
		t.Errorf("recorded response %s = %q, want %q", payloadboxRecordIDHeader, gotHeader, got.ID)
	}
}

func TestRecord_GETUnknownPathIsCaptured(t *testing.T) {
	handler, s := newTestHandler(t, defaultSettings())

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/aaaaz", http.NoBody)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202", rec.Code)
	}

	if got := rec.Body.String(); got != "Request logged GET /aaaaz\n" {
		t.Fatalf("body = %q, want recorder response", got)
	}

	records := s.store.List()
	if len(records) != 1 {
		t.Fatalf("store has %d records, want 1", len(records))
	}

	if records[0].Request.Path != "/aaaaz" {
		t.Errorf("path = %q, want /aaaaz", records[0].Request.Path)
	}
}

func TestRecord_UIRootBypassesCapture(t *testing.T) {
	handler, s := newTestHandler(t, defaultSettings())

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", http.NoBody)
	handler.ServeHTTP(rec, req)

	if rec.Code == http.StatusAccepted {
		t.Fatalf("status = %d, want UI handler response", rec.Code)
	}

	if records := s.store.List(); len(records) != 0 {
		t.Fatalf("store has %d records, want 0", len(records))
	}
}

func TestRecord_BodyTooLargeStillCaptured(t *testing.T) {
	settings := defaultSettings()
	settings.MaxBodySizeBytes = 10
	handler, s := newTestHandler(t, settings)

	big := strings.Repeat("x", 100)
	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/big", strings.NewReader(big))
	req.Header.Set("Content-Type", "text/plain")
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}

	records := s.store.List()
	if len(records) != 1 {
		t.Fatalf("store has %d records, want 1 (attempt should still be recorded)", len(records))
	}

	if records[0].Response.StatusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("recorded status = %d, want 413", records[0].Response.StatusCode)
	}

	if headerID := rec.Header().Get(payloadboxRecordIDHeader); headerID != records[0].ID {
		t.Errorf("%s = %q, want stored record ID %q", payloadboxRecordIDHeader, headerID, records[0].ID)
	}
}

func TestListRecords_ReturnsJSONArray(t *testing.T) {
	handler, s := newTestHandler(t, defaultSettings())

	// Seed two captures via the public API (also exercises Record).
	for _, p := range []string{"/a", "/b"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, p, strings.NewReader("x"))
		req.Header.Set("Content-Type", "text/plain")
		handler.ServeHTTP(rec, req)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/history", http.NoBody)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var records []Record

	if err := json.NewDecoder(rec.Body).Decode(&records); err != nil {
		t.Fatal(err)
	}

	if len(records) != 2 {
		t.Fatalf("got %d records, want 2", len(records))
	}

	if len(s.store.List()) != len(records) {
		t.Errorf("server store has %d, response has %d", len(s.store.List()), len(records))
	}
}

func TestGetRecord_ByID(t *testing.T) {
	handler, s := newTestHandler(t, defaultSettings())

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/x", strings.NewReader("hi"))
	req.Header.Set("Content-Type", "text/plain")
	handler.ServeHTTP(rec, req)

	stored := s.store.List()
	if len(stored) != 1 {
		t.Fatalf("store has %d records, want 1", len(stored))
	}

	id := stored[0].ID
	if id == "" {
		t.Fatal("record has empty ID")
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/history/"+id, http.NoBody)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var got Record
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if got.ID != id {
		t.Errorf("returned ID = %q, want %q", got.ID, id)
	}

	if got.Request.Path != "/x" {
		t.Errorf("returned path = %q, want /x", got.Request.Path)
	}
}

func TestGetRecord_UnknownIDReturns404(t *testing.T) {
	handler, _ := newTestHandler(t, defaultSettings())

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/history/does-not-exist", http.NoBody)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestClearRecords(t *testing.T) {
	handler, s := newTestHandler(t, defaultSettings())

	// Seed one.
	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/x", strings.NewReader("x"))
	req.Header.Set("Content-Type", "text/plain")
	handler.ServeHTTP(rec, req)

	if len(s.store.List()) == 0 {
		t.Fatal("expected record before DELETE")
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/api/v1/history", http.NoBody)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}

	if n := len(s.store.List()); n != 0 {
		t.Fatalf("store has %d records after DELETE, want 0", n)
	}
}

func TestHealth(t *testing.T) {
	handler, _ := newTestHandler(t, defaultSettings())

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/healthz", http.NoBody)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	if rec.Body.String() != "OK" {
		t.Errorf("body = %q, want %q", rec.Body.String(), "OK")
	}
}

func TestVersion(t *testing.T) {
	handler, _ := newTestHandler(t, defaultSettings())

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/version", http.NoBody)
	handler.ServeHTTP(rec, req)

	var out map[string]string

	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	for _, k := range []string{"version", "build_sha", "build_time"} {
		if _, ok := out[k]; !ok {
			t.Errorf("response missing key %q; got %v", k, out)
		}
	}
}

func TestSettings(t *testing.T) {
	s := defaultSettings()
	s.Address = ":9999"
	s.MaxBodySizeBytes = 2048
	handler, _ := newTestHandler(t, s)

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/settings", http.NoBody)
	handler.ServeHTTP(rec, req)

	var out map[string]any

	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	if got := out["address"]; got != ":9999" {
		t.Errorf("address = %v, want :9999", got)
	}

	if got, want := out["max_body_size_bytes"], float64(2048); got != want {
		t.Errorf("max_body_size_bytes = %v, want %v", got, want)
	}
}

func TestEvents_EmitsRecordFrameOnPublish(t *testing.T) {
	ts, s := newTestServer(t, defaultSettings())

	// Subscribe first (start the SSE connection), then publish.
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, ts.URL+"/api/v1/events", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		t.Fatalf("Content-Type = %q, want text/event-stream", ct)
	}

	// Give the subscriber a moment to register with the hub before publishing.
	// (httptest dispatches in a goroutine; the handler does Subscribe() on first line.)
	time.Sleep(50 * time.Millisecond)

	// Publish by recording a request through the public API.
	go func() {
		pubReq, err := http.NewRequestWithContext(t.Context(), http.MethodPost, ts.URL+"/hit", strings.NewReader("x"))
		if err != nil {
			return
		}

		pubReq.Header.Set("Content-Type", "text/plain")

		pubResp, err := http.DefaultClient.Do(pubReq)
		if err == nil {
			_ = pubResp.Body.Close()
		}
	}()

	reader := bufio.NewReader(resp.Body)
	done := make(chan struct{})

	var gotEvent, gotData string

	go func() {
		defer close(done)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}

			line = strings.TrimRight(line, "\r\n")
			if v, ok := strings.CutPrefix(line, "event: "); ok {
				gotEvent = v
			}

			if v, ok := strings.CutPrefix(line, "data: "); ok {
				gotData = v
				return
			}
		}
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("no SSE frame within 2s")
	}

	if gotEvent != "record" {
		t.Errorf("event = %q, want record", gotEvent)
	}

	var rec Record

	if err := json.Unmarshal([]byte(gotData), &rec); err != nil {
		t.Fatalf("invalid JSON in SSE data: %v (raw=%q)", err, gotData)
	}

	if rec.Request.Path != "/hit" {
		t.Errorf("recorded path = %q, want /hit", rec.Request.Path)
	}

	// Sanity: store observed the record too.
	if n := len(s.store.List()); n == 0 {
		t.Error("store empty after publish; expected at least 1")
	}
}
