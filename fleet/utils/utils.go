package utils

type comparable interface {
	int | ~string
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

func RemoveByIdx[T any](arr []T, idx int) []T {
	return append(arr[:idx], arr[idx+1:]...)
}

func FindAndRemove[T comparable](arr []T, value T) []T {
	out := []T{}
	for _, v := range arr {
		if v != value {
			out = append(out, v)
		}
	}
	return out
}

func AddUnique[T comparable](arr []T, value T) []T {
	out := FindAndRemove(arr, value)
	out = append(out, value)
	return out
}
