package util

// Default checks if the given value is empty, and if so returns the defaultValue
func Default[T comparable](value, defaultValue T) T {
	var x T
	if value == x {
		return defaultValue
	}
	return value
}

func Ptr[T any](value T) *T { return &value }
