package validation

import (
	"regexp"
)

func ValidateName(name string) bool {
	b, _ := regexp.MatchString("^[a-z0-9\\-]*$", name)
	return b
}
