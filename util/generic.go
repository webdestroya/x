package util

// Bring back one liner ternary
//
//	bool ? trueVal : falseVal
func Ternary[T any](expr bool, trueVal T, falseVal T) T {
	if expr {
		return trueVal
	}
	return falseVal
}

// returns the first value of a multi-return statement
func FirstParam[T any, U any](fp T, _ ...U) T {
	return fp
}

// Return or panic
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
