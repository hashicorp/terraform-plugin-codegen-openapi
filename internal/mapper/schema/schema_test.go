package schema_test

// TODO: holding here for safe-keeping
func pointer[T any](value T) *T {
	return &value
}
