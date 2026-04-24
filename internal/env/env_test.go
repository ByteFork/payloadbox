package env

import (
	"testing"
)

func TestString(t *testing.T) {
	t.Setenv("X", "hello")

	if got := String("X", "fallback"); got != "hello" {
		t.Errorf("String set = %q, want hello", got)
	}

	t.Setenv("X", "")

	if got := String("X", "fallback"); got != "fallback" {
		t.Errorf("String empty = %q, want fallback", got)
	}
}

func TestBool(t *testing.T) {
	cases := []struct {
		in   string
		want bool
		fb   bool
	}{
		{"true", true, false},
		{"false", false, true},
		{"1", true, false},
		{"0", false, true},
		{"garbage", true, true},
		{"garbage", false, false},
		{"", true, true},
	}
	for _, c := range cases {
		t.Setenv("X", c.in)

		if got := Bool("X", c.fb); got != c.want {
			t.Errorf("Bool(%q, fb=%v) = %v, want %v", c.in, c.fb, got, c.want)
		}
	}
}

func TestInt(t *testing.T) {
	cases := []struct {
		in   string
		fb   int
		want int
	}{
		{"42", 0, 42},
		{"-7", 0, -7},
		{"", 99, 99},
		{"notanum", 99, 99},
	}
	for _, c := range cases {
		t.Setenv("X", c.in)

		if got := Int("X", c.fb); got != c.want {
			t.Errorf("Int(%q, fb=%d) = %d, want %d", c.in, c.fb, got, c.want)
		}
	}
}

func TestInt64(t *testing.T) {
	cases := []struct {
		in   string
		fb   int64
		want int64
	}{
		{"9999999999", 0, 9999999999},
		{"", 1, 1},
		{"notanum", 1, 1},
	}
	for _, c := range cases {
		t.Setenv("X", c.in)

		if got := Int64("X", c.fb); got != c.want {
			t.Errorf("Int64(%q, fb=%d) = %d, want %d", c.in, c.fb, got, c.want)
		}
	}
}
