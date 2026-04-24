// Package logutil provides small helpers for structured logging.
package logutil

import (
	"log/slog"
	"strings"
)

// Level maps a string to slog.Level. Unknown values default to Info.
// Accepts (case-insensitive): "debug", "info", "warn", "warning", "error".
func Level(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// SafeField strips C0 control characters and DEL from s, returning a string
// safe to include as a structured log field value. Use it to wrap any field
// derived from HTTP input (paths, headers, query strings, error messages
// carrying user data) before passing to slog, to prevent log forging via
// embedded newlines or terminal escape sequences.
func SafeField(s string) string {
	return strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}

		return r
	}, s)
}
