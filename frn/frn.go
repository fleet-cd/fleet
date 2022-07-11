package frn

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type String = string

type FRN struct {
	Resource  string `json:"resource" yaml:"frn"`
	Namespace string `json:"namespace" yaml:"frn"`
	Locator   string `json:"locator" yaml:"frn"`
}

func (frn FRN) String() String {
	return String(fmt.Sprintf("frn:%s:%s:%s", frn.Resource, frn.Namespace, frn.Locator))
}

func Generate(resource, namespace string) FRN {
	return FRN{
		Resource:  resource,
		Namespace: namespace,
		Locator:   uuid.New().String(),
	}
}

func Parse(frn String) (FRN, error) {
	slice := strings.Split(frn, ":")
	if slice[0] != "frn" || len(slice) != 4 {
		return FRN{}, fmt.Errorf("string \"%s\" is not a proper FRN", frn)
	}
	return FRN{
		Resource:  slice[1],
		Namespace: slice[2],
		Locator:   slice[3],
	}, nil
}
