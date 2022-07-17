package panicker

type Handler[T any] struct {
	err   error
	value T
}

func CheckAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckAndPanicFunc(err error, fcn func(err error) error) {
	if err != nil {
		panic(fcn(err))
	}
}

func (h Handler[T]) GetOrPanicFunc(fcn func(err error) error) T {
	if h.err != nil {
		panic(fcn(h.err))
	}
	return h.value
}

func (h Handler[T]) GetOrPanic() T {
	if h.err != nil {
		panic(h.err)
	}
	return h.value
}

func AndPanic[T any](val T, err error) Handler[T] {
	return Handler[T]{
		err:   err,
		value: val,
	}
}

func OnFalse(val bool, err error) {
	if val == false {
		panic(err)
	}
}

func Panicker[T any](val T, err error) Handler[T] {
	return Handler[T]{
		err:   err,
		value: val,
	}
}
