package utils

func CheckAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}
