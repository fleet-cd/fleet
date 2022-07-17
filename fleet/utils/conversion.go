package utils

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func ConvertList[I any, O any](inputs []I, conversionFcn func(input I) O) []O {
	out := []O{}
	for _, i := range inputs {
		out = append(out, conversionFcn(i))
	}
	return out
}

func GetSortMap(sort string) bson.D {
	sortMap := bson.D{}
	if sort != "" {
		sortDir := -1 // descending
		if sort[0] == '+' {
			sortDir = 1
			sort = strings.TrimPrefix(sort, "+")
		} else if sort[0] == '-' {
			sort = strings.TrimPrefix(sort, "-")
		}
		sortMap = append(sortMap, bson.E{Key: sort, Value: sortDir})
	}
	return sortMap
}
