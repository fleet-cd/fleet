package utils

type comparable interface {
	int | string
}

func OrDefault[T any](pointer *T, value T) T {
	if pointer != nil {
		return *pointer
	}
	return value
}

func Contains[T comparable](arr []T, value T) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
