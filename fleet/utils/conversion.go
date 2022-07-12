package utils

func ConvertList[I any, O any](inputs []I, conversionFcn func(input I) O) []O {
	out := []O{}
	for _, i := range inputs {
		out = append(out, conversionFcn(i))
	}
	return out
}
