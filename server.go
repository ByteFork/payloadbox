package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/ByteFork/payloadbox/internal/hub"
	"github.com/ByteFork/payloadbox/internal/store"
)

const payloadboxRecordIDHeader = "X-Payloadbox-Record-ID"

type (
	Server struct {
		settings ServerSettings
		store    *store.Store[Record]
		hub      *hub.Hub[Record]
	}

	RequestData struct {
		Method     string              `json:"method"`
		Path       string              `json:"path"`
		Query      string              `json:"query,omitempty"`
		Headers    map[string][]string `json:"headers"`
		RemoteAddr string              `json:"remote_addr"`
		UserAgent  string              `json:"user_agent"`
		Body       string              `json:"body,omitempty"`
	}

	ResponseData struct {
		StatusCode int                 `json:"status_code"`
		StatusText string              `json:"status_text"`
		Headers    map[string][]string `json:"headers,omitempty"`
		Body       string              `json:"body,omitempty"`
		Size       int64               `json:"size_in_bytes"`
	}

	Record struct {
		ID        string       `json:"id"`
		CreatedAt time.Time    `json:"created_at"`
		Duration  int64        `json:"duration_ns"`
		Request   RequestData  `json:"request"`
		Response  ResponseData `json:"response"`
	}
)

func NewServer(settings ServerSettings) *Server {
	return &Server{
		settings: settings,
		store:    store.NewStore[Record](settings.MaxRecordsToStore),
		hub:      hub.NewHub[Record](),
	}
}

func NewRecord(req *http.Request, reqBody string, status int, respHeaders http.Header, respBody string, duration, size int64) Record {
	if req == nil {
		return Record{}
	}

	var respHeadersMap map[string][]string
	if respHeaders != nil {
		respHeadersMap = respHeaders.Clone()
	}

	return Record{
		ID:        newRecordID(),
		CreatedAt: time.Now().UTC(),
		Duration:  duration,
		Request: RequestData{
			Method:     req.Method,
			Path:       req.URL.Path,
			Query:      req.URL.RawQuery,
			Headers:    req.Header.Clone(),
			RemoteAddr: req.RemoteAddr,
			UserAgent:  req.UserAgent(),
			Body:       reqBody,
		},
		Response: ResponseData{
			StatusCode: status,
			StatusText: http.StatusText(status),
			Headers:    respHeadersMap,
			Body:       respBody,
			Size:       size,
		},
	}
}

// newRecordID returns a stable, time-ordered identifier for a captured request.
// UUIDv7 sorts as bytes by creation time, which makes it a natural primary key
// for the in-memory store and any future persistence. A random UUIDv4 is used
// as a fallback if v7 generation fails (extremely rare; documented in
// google/uuid as a clock/entropy edge case).
func newRecordID() string {
	if id, err := uuid.NewV7(); err == nil {
		return id.String()
	}

	return uuid.NewString()
}
