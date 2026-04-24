// Package env reads environment variables with fallback defaults.
// A variable that is unset, empty, or fails to parse yields the fallback.
package env

import (
	"os"
	"strconv"
)

// String returns the value of key, or fallback if unset or empty.
func String(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

// Bool parses key as a bool, returning fallback on empty or parse error.
func Bool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}

	return b
}

// Int parses key as an int, returning fallback on empty or parse error.
func Int(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}

	return n
}

// Int64 parses key as an int64, returning fallback on empty or parse error.
func Int64(key string, fallback int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return fallback
	}

	return n
}
