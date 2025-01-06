package gocrew

func prepend[T any](array []T, value T) []T {
	array = append(array, *new(T))
	copy(array[1:], array)
	array[0] = value
	return array
}

// P wraps the value to a pointer
func P[T any](v T) *T {
	return &v
}
