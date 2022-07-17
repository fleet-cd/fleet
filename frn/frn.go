package frn

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
)

type String = string

type FRN struct {
	Resource  string `json:"resource" yaml:"frn"`
	Namespace string `json:"namespace" yaml:"frn"`
	Locator   string `json:"locator" yaml:"frn"`
}

func Namespace(f common.Frn) string {
	split := strings.Split(string(f), ":")
	return split[2]
}

func Resource(f common.Frn) string {
	split := strings.Split(string(f), ":")
	return split[3]
}

func (frn FRN) String() common.Frn {
	return common.Frn(fmt.Sprintf("frn:%s:%s:%s", frn.Resource, frn.Namespace, frn.Locator))
}

func Generate(resource, namespace string) FRN {
	return FRN{
		Resource:  resource,
		Namespace: namespace,
		Locator:   uuid.New().String(),
	}
}

func GenerateActual[T ~string](resource, namespace string) T {
	return T(FRN{
		Resource:  resource,
		Namespace: namespace,
		Locator:   uuid.New().String(),
	}.String())
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
