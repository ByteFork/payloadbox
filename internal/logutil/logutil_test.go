package logutil

import (
	"log/slog"
	"testing"
)

func TestLevel(t *testing.T) {
	cases := []struct {
		in   string
		want slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"INFO", slog.LevelInfo},
		{"", slog.LevelInfo},
		{"garbage", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"warning", slog.LevelWarn},
		{"WARNING", slog.LevelWarn},
		{"error", slog.LevelError},
		{"ERROR", slog.LevelError},
	}

	for _, c := range cases {
		if got := Level(c.in); got != c.want {
			t.Errorf("Level(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestSafeField(t *testing.T) {
	cases := []struct {
		name     string
		in, want string
	}{
		{"empty", "", ""},
		{"plain ascii", "hello", "hello"},
		{"typical path", "/api/v1/users", "/api/v1/users"},
		{"strips newline", "foo\nbar", "foobar"},
		{"strips carriage return", "foo\rbar", "foobar"},
		{"strips tab", "foo\tbar", "foobar"},
		{"strips null byte", "x\x00y", "xy"},
		{"strips DEL", "safe\x7fdel", "safedel"},
		{"strips ANSI escape", "red\x1b[31m", "red[31m"},
		{"keeps unicode", "unicode-é-safe", "unicode-é-safe"},
		{"log forging attempt", "/foo\n[ERROR] admin=true", "/foo[ERROR] admin=true"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := SafeField(c.in); got != c.want {
				t.Errorf("SafeField(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
