package auth

import (
	"fmt"

	"github.com/tgs266/fleet/fleet/utils"
	"github.com/tgs266/rest-gen/runtime/authentication"
)

type Action = string

const (
	ACTION_VIEW   Action = "VIEW"
	ACTION_EDIT   Action = "EDIT"
	ACTION_CREATE Action = "CREATE"
	ACTION_DELETE Action = "DELETE"
)

func CanI(authToken authentication.Token, resource string, namespace string, action Action) (bool, error) {
	jwt, err := DecodeToken(authToken)
	if err != nil {
		return false, err
	}
	user, err := GetUserByFrn(jwt.UserFrn)
	if err != nil {
		return false, err
	}
	if utils.Contains(user.Groups, "admins") {
		return true, nil
	}
	return false, fmt.Errorf("bad permissions")
}

func CanIFrn(authToken authentication.Token, resourceFrn string, action Action) (bool, error) {
	jwt, err := DecodeToken(authToken)
	if err != nil {
		return false, err
	}
	user, err := GetUserByFrn(jwt.UserFrn)
	if err != nil {
		return false, err
	}
	if utils.Contains(user.Groups, "admins") {
		return true, nil
	}
	return false, fmt.Errorf("bad permissions")
}
