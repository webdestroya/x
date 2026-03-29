package util

type stringish interface {
	string | *string
}

// Returns the first non-empty string
func CoalesceStrings(strs ...string) string {
	for _, str := range strs {
		if str != "" {
			return str
		}
	}
	return ""
}

// Returns the first non nil and non empty string
func Coalesce[T stringish](vals ...T) T {
	for _, val := range vals {
		switch v := any(val).(type) {
		case string:
			if v != "" {
				return val
			}
		case *string:
			if v != nil && *v != "" {
				return val
			}
		}
	}

	var zero T
	return zero
}

// Returns the first non nil value
func CoalescePointers[T any](vals ...*T) *T {
	for _, v := range vals {
		if v != nil {
			return v
		}
	}
	return nil
}

// Returns the first non-nil element that the provided function returns true for
func CoalesceWithFunc[T any](checkFunc func(T) bool, values ...T) (T, bool) {
	for _, value := range values {
		if checkFunc(value) {
			return value, true
		}
	}
	var zero T
	return zero, false
}
