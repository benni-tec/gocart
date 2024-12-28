package gotrac

// P wraps the value to a pointer
func P[T any](v T) *T {
	return &v
}
