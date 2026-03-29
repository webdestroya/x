package util_test

import (
	"errors"
	"testing"

	"github.com/webdestroya/x/util"
)

func TestTernary(t *testing.T) {
	t.Run("int true", func(t *testing.T) {
		got := util.Ternary(true, 1, 2)
		if got != 1 {
			t.Errorf("Ternary(true, 1, 2) = %d, want 1", got)
		}
	})

	t.Run("int false", func(t *testing.T) {
		got := util.Ternary(false, 1, 2)
		if got != 2 {
			t.Errorf("Ternary(false, 1, 2) = %d, want 2", got)
		}
	})

	t.Run("string true", func(t *testing.T) {
		got := util.Ternary(true, "yes", "no")
		if got != "yes" {
			t.Errorf("Ternary(true, yes, no) = %q, want %q", got, "yes")
		}
	})

	t.Run("string false", func(t *testing.T) {
		got := util.Ternary(false, "yes", "no")
		if got != "no" {
			t.Errorf("Ternary(false, yes, no) = %q, want %q", got, "no")
		}
	})
}

func TestFirstParam(t *testing.T) {
	t.Run("int with single trailing arg", func(t *testing.T) {
		got := util.FirstParam(1, "ignored")
		if got != 1 {
			t.Errorf("FirstParam(1, ...) = %d, want 1", got)
		}
	})

	t.Run("string with error", func(t *testing.T) {
		got := util.FirstParam("ok", errors.New("err"))
		if got != "ok" {
			t.Errorf("FirstParam(ok, ...) = %q, want %q", got, "ok")
		}
	})

	t.Run("int with multiple trailing args", func(t *testing.T) {
		got := util.FirstParam(42, "a", "b")
		if got != 42 {
			t.Errorf("FirstParam(42, ...) = %d, want 42", got)
		}
	})
}

func TestMust(t *testing.T) {
	t.Run("returns value on nil error", func(t *testing.T) {
		got := util.Must(42, nil)
		if got != 42 {
			t.Errorf("Must(42, nil) = %d, want 42", got)
		}
	})

	t.Run("returns string on nil error", func(t *testing.T) {
		got := util.Must("hello", nil)
		if got != "hello" {
			t.Errorf("Must(hello, nil) = %q, want %q", got, "hello")
		}
	})

	t.Run("panics on error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Must() did not panic on non-nil error")
			}
		}()
		util.Must(0, errors.New("boom"))
	})
}
