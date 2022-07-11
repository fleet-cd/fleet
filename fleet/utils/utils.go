package utils

func OrDefault[T any](pointer *T, value T) T {
	if pointer != nil {
		return *pointer
	}
	return value
}
