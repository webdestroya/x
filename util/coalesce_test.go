package util_test

import (
	"testing"

	"github.com/webdestroya/x/util"
)

func TestCoalesceStrings(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want string
	}{
		{"first non-empty", []string{"a", "b"}, "a"},
		{"skips leading empties", []string{"", "", "c"}, "c"},
		{"all empty", []string{"", ""}, ""},
		{"no args", nil, ""},
		{"single non-empty", []string{"x"}, "x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.CoalesceStrings(tt.in...)
			if got != tt.want {
				t.Errorf("CoalesceStrings(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestCoalesce_String(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want string
	}{
		{"first non-empty", []string{"a", "b"}, "a"},
		{"skips leading empties", []string{"", "", "c"}, "c"},
		{"all empty", []string{"", ""}, ""},
		{"no args", nil, ""},
		{"single non-empty", []string{"x"}, "x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.Coalesce(tt.in...)
			if got != tt.want {
				t.Errorf("Coalesce(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestCoalesce_StringPtr(t *testing.T) {
	tests := []struct {
		name string
		in   []*string
		want *string
	}{
		{"first valid", []*string{new("a"), new("b")}, new("a")},
		{"skips nil", []*string{nil, new("b")}, new("b")},
		{"skips empty ptr", []*string{new(""), new("c")}, new("c")},
		{"skips nil and empty", []*string{nil, new(""), new("d")}, new("d")},
		{"all nil", []*string{nil, nil}, nil},
		{"all empty ptrs", []*string{new(""), new("")}, nil},
		{"no args", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.Coalesce(tt.in...)
			if tt.want == nil {
				if got != nil {
					t.Errorf("Coalesce() = %q, want nil", *got)
				}
				return
			}
			if got == nil {
				t.Errorf("Coalesce() = nil, want %q", *tt.want)
				return
			}
			if *got != *tt.want {
				t.Errorf("Coalesce() = %q, want %q", *got, *tt.want)
			}
		})
	}
}

func TestCoalescePointers(t *testing.T) {
	tests := []struct {
		name string
		in   []*int
		want *int
	}{
		{"first non-nil", []*int{new(1), new(2)}, new(1)},
		{"skips leading nils", []*int{nil, nil, new(3)}, new(3)},
		{"all nil", []*int{nil, nil}, nil},
		{"no args", nil, nil},
		{"single non-nil", []*int{new(42)}, new(42)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.CoalescePointers(tt.in...)
			if tt.want == nil {
				if got != nil {
					t.Errorf("CoalescePointers() = %d, want nil", *got)
				}
				return
			}
			if got == nil {
				t.Errorf("CoalescePointers() = nil, want %d", *tt.want)
				return
			}
			if *got != *tt.want {
				t.Errorf("CoalescePointers() = %d, want %d", *got, *tt.want)
			}
		})
	}
}

func TestCoalesceWithFunc(t *testing.T) {
	positive := func(n int) bool { return n > 0 }

	t.Run("returns first match", func(t *testing.T) {
		got, ok := util.CoalesceWithFunc(positive, -1, 0, 5, 10)
		if !ok {
			t.Fatal("CoalesceWithFunc() ok = false, want true")
		}
		if got != 5 {
			t.Errorf("CoalesceWithFunc() = %d, want 5", got)
		}
	})

	t.Run("first value matches", func(t *testing.T) {
		got, ok := util.CoalesceWithFunc(positive, 3, -1, 0)
		if !ok {
			t.Fatal("CoalesceWithFunc() ok = false, want true")
		}
		if got != 3 {
			t.Errorf("CoalesceWithFunc() = %d, want 3", got)
		}
	})

	t.Run("no match", func(t *testing.T) {
		got, ok := util.CoalesceWithFunc(positive, -1, -2, 0)
		if ok {
			t.Fatal("CoalesceWithFunc() ok = true, want false")
		}
		if got != 0 {
			t.Errorf("CoalesceWithFunc() = %d, want 0", got)
		}
	})

	t.Run("no args", func(t *testing.T) {
		got, ok := util.CoalesceWithFunc(positive)
		if ok {
			t.Fatal("CoalesceWithFunc() ok = true, want false")
		}
		if got != 0 {
			t.Errorf("CoalesceWithFunc() = %d, want 0", got)
		}
	})
}
