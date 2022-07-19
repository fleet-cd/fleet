package utils

import "fmt"

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

func Values[T any](obj map[string]T) []T {
	vals := []T{}
	for _, v := range obj {
		vals = append(vals, v)
	}
	return vals
}

func Find[T any](arr []T, predicate func(i T) bool) (T, error) {
	for _, v := range arr {
		if predicate(v) {
			return v, nil
		}
	}
	return arr[0], fmt.Errorf("not found")
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
