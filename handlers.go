package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ByteFork/payloadbox/internal/logutil"
)

// sseKeepaliveInterval is how often the SSE handler emits a heartbeat comment
// to keep intermediaries (proxies, load balancers) from dropping idle streams.
const sseKeepaliveInterval = 30 * time.Second

var (
	sseRecordPrefix = []byte("event: record\ndata: ")
	sseFrameEnd     = []byte("\n\n")
	sseHeartbeat    = []byte(": ping\n\n")
)

func (s *Server) Record(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if r.Method == http.MethodGet && !strings.HasPrefix(r.URL.Path, "/api/") {
		s.UI(w, r)
		return
	}

	if s.settings.LogRequests {
		slog.Info("incoming request",
			"method", logutil.SafeField(r.Method),
			"path", logutil.SafeField(r.URL.Path),
		)
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			slog.Error("close request body", "error", logutil.SafeField(err.Error()))
		}
	}()

	r.Body = http.MaxBytesReader(w, r.Body, s.settings.MaxBodySizeBytes)
	body, err := io.ReadAll(r.Body)
	status := http.StatusAccepted

	if err != nil {
		slog.Error("failed to read request body",
			"error", logutil.SafeField(err.Error()),
			"path", logutil.SafeField(r.URL.Path),
		)

		status = http.StatusRequestEntityTooLarge
		respBody := "Request Body Too Large"

		// Record the attempt even if it failed due to size
		duration := time.Since(start).Nanoseconds()
		record := NewRecord(r, "", status, nil, respBody, duration, int64(len(respBody)))
		s.store.Add(record)
		s.hub.Publish(record)

		http.Error(w, respBody, status)

		return
	}

	respBody := fmt.Sprintf("Request logged %s %s\n", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	duration := time.Since(start).Nanoseconds()
	record := NewRecord(r, string(body), status, w.Header(), respBody, duration, int64(len(respBody)))
	s.store.Add(record)
	s.hub.Publish(record)

	w.WriteHeader(status)

	if _, err := fmt.Fprint(w, respBody); err != nil {
		slog.Error("failed to write response body",
			"error", logutil.SafeField(err.Error()),
			"path", logutil.SafeField(r.URL.Path),
		)
	}

	if s.settings.LogRequests {
		slog.Info("incoming request recorded",
			"path", logutil.SafeField(r.URL.Path),
			"method", logutil.SafeField(r.Method),
		)
	}
}

func (s *Server) ListRecords(w http.ResponseWriter, _ *http.Request) {
	records := s.store.List()

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := enc.Encode(records); err != nil {
		http.Error(w, "Failed to encode records", http.StatusInternalServerError)

		return
	}
}

func (s *Server) ClearRecords(w http.ResponseWriter, _ *http.Request) {
	s.store.Clear()
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) Events(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Opt this long-lived stream out of the server's WriteTimeout.
	rc := http.NewResponseController(w)
	_ = rc.SetWriteDeadline(time.Time{})

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	ch := s.hub.Subscribe()
	defer s.hub.Unsubscribe(ch)

	keepalive := time.NewTicker(sseKeepaliveInterval)
	defer keepalive.Stop()

	for {
		select {
		case rec := <-ch:
			b, err := json.Marshal(rec)
			if err != nil {
				slog.Error("SSE marshal record", "error", logutil.SafeField(err.Error()))
				continue
			}

			frame := make([]byte, 0, len(sseRecordPrefix)+len(b)+len(sseFrameEnd))
			frame = append(frame, sseRecordPrefix...)
			frame = append(frame, b...)
			frame = append(frame, sseFrameEnd...)
			_, _ = w.Write(frame)

			flusher.Flush()
		case <-keepalive.C:
			_, _ = w.Write(sseHeartbeat)

			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	if s.settings.LogRequests {
		slog.Info("incoming request",
			"path", logutil.SafeField(r.URL.Path),
			"method", logutil.SafeField(r.Method),
		)
	}

	if _, err := fmt.Fprint(w, "OK"); err != nil {
		slog.Error("failed to write response body",
			"error", logutil.SafeField(err.Error()),
			"path", logutil.SafeField(r.URL.Path),
		)
	}
}

func (s *Server) Version(w http.ResponseWriter, r *http.Request) {
	if s.settings.LogRequests {
		slog.Info("incoming request",
			"path", logutil.SafeField(r.URL.Path),
			"method", logutil.SafeField(r.Method),
		)
	}

	if err := json.NewEncoder(w).Encode(map[string]string{
		"version":    Version,
		"build_sha":  BuildSha,
		"build_time": BuildTime,
	}); err != nil {
		slog.Error("failed to write response body",
			"error", logutil.SafeField(err.Error()),
			"path", logutil.SafeField(r.URL.Path),
		)
	}
}

func (s *Server) Settings(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(map[string]any{
		"address":              s.settings.Address,
		"max_body_size_bytes":  s.settings.MaxBodySizeBytes,
		"max_records_to_store": s.settings.MaxRecordsToStore,
		"log_requests":         s.settings.LogRequests,
		"log_level":            s.settings.LogLevel,
	}); err != nil {
		slog.Error("failed to write settings", "error", logutil.SafeField(err.Error()))
	}
}
