package main

import (
	"testing"
)

func TestNewSettings_Defaults(t *testing.T) {
	// t.Setenv would leak if someone else set these, so clear them first.
	for _, k := range []string{
		"LISTEN_ADDRESS", "MAX_BODY_SIZE_BYTES", "MAX_RECORDS_TO_STORE",
		"LOG_HTTP_REQUESTS", "LOG_LEVEL",
	} {
		t.Setenv(k, "")
	}

	s := NewSettings()
	if s.Address != ":8080" {
		t.Errorf("Address = %q, want %q", s.Address, ":8080")
	}

	if s.MaxBodySizeBytes != 1024 {
		t.Errorf("MaxBodySizeBytes = %d, want 1024", s.MaxBodySizeBytes)
	}

	if s.MaxRecordsToStore != 200 {
		t.Errorf("MaxRecordsToStore = %d, want 200", s.MaxRecordsToStore)
	}

	if !s.LogRequests {
		t.Error("LogRequests = false, want true")
	}

	if s.LogLevel != "info" {
		t.Errorf("LogLevel = %q, want info", s.LogLevel)
	}
}

func TestNewSettings_EnvOverrides(t *testing.T) {
	t.Setenv("LISTEN_ADDRESS", ":9000")
	t.Setenv("MAX_BODY_SIZE_BYTES", "5000")
	t.Setenv("MAX_RECORDS_TO_STORE", "42")
	t.Setenv("LOG_HTTP_REQUESTS", "false")
	t.Setenv("LOG_LEVEL", "debug")

	s := NewSettings()
	if s.Address != ":9000" {
		t.Errorf("Address = %q, want :9000", s.Address)
	}

	if s.MaxBodySizeBytes != 5000 {
		t.Errorf("MaxBodySizeBytes = %d, want 5000", s.MaxBodySizeBytes)
	}

	if s.MaxRecordsToStore != 42 {
		t.Errorf("MaxRecordsToStore = %d, want 42", s.MaxRecordsToStore)
	}

	if s.LogRequests {
		t.Error("LogRequests = true, want false")
	}

	if s.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want debug", s.LogLevel)
	}
}

func TestNewSettings_InvalidEnvFallsBackToDefault(t *testing.T) {
	t.Setenv("MAX_BODY_SIZE_BYTES", "notanumber")
	t.Setenv("MAX_RECORDS_TO_STORE", "xxx")
	t.Setenv("LOG_HTTP_REQUESTS", "maybe")

	s := NewSettings()
	if s.MaxBodySizeBytes != 1024 {
		t.Errorf("MaxBodySizeBytes = %d, want default 1024", s.MaxBodySizeBytes)
	}

	if s.MaxRecordsToStore != 200 {
		t.Errorf("MaxRecordsToStore = %d, want default 200", s.MaxRecordsToStore)
	}

	if !s.LogRequests {
		t.Error("LogRequests = false, want default true")
	}
}
