package auth

import (
	"fmt"

	"github.com/tgs266/fleet/fleet/utils"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"github.com/tgs266/rest-gen/runtime/authentication"
	"github.com/tgs266/rest-gen/runtime/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type Action = string

const (
	ACTION_VIEW   Action = "VIEW"
	ACTION_EDIT   Action = "EDIT"
	ACTION_CREATE Action = "CREATE"
	ACTION_DELETE Action = "DELETE"
)

func IsAuth(authToken authentication.Token) (entities.UserEntity, error) {
	jwt, err := DecodeToken(authToken)
	if err != nil {
		return entities.UserEntity{}, err
	}
	user, err := GetUserByFrn(jwt.UserFrn)
	if err != nil {
		return entities.UserEntity{}, errors.NewUnauthorized(err)
	}
	return user, nil
}

func WhatCanIList(authToken authentication.Token, resource string) (bson.M, error) {
	user, err := IsAuth(authToken)
	if err != nil {
		return nil, err
	}
	if utils.Contains(user.Groups, "admins") {
		return bson.M{}, nil
	}

	groups, err := GetGroups(user.Groups)
	if err != nil {
		return nil, err
	}
	groupString := []common.PermissionFrn{}
	for _, g := range groups {
		groupString = append(groupString, g.Permissions...)
	}

	perms, err := GetPermissions(groupString)
	accessNamespaces := []string{}
	for _, p := range perms {
		if p.ResourceType == resource && utils.Contains(p.Actions, "view") {
			if p.Namespace == "*" {
				return bson.M{}, nil
			}
			accessNamespaces = append(accessNamespaces, p.Namespace)
		}
	}

	return bson.M{"namespace": bson.M{"$in": accessNamespaces}}, nil
}

func CanI(authToken authentication.Token, resource string, namespace string, action Action) (bool, error) {
	user, err := IsAuth(authToken)
	if err != nil {
		return false, err
	}
	if utils.Contains(user.Groups, "admins") {
		return true, nil
	}
	return false, fmt.Errorf("bad permissions")
}

func IsAdmin(authToken authentication.Token) (bool, error) {
	user, err := IsAuth(authToken)
	if err != nil {
		return false, err
	}
	if utils.Contains(user.Groups, "admins") {
		return true, nil
	}
	return false, errors.NewForbidden(nil)
}

func CanIFrn(authToken authentication.Token, resourceFrn string, action Action) (bool, error) {
	user, err := IsAuth(authToken)
	if err != nil {
		return false, err
	}
	if utils.Contains(user.Groups, "admins") {
		return true, nil
	}
	return false, fmt.Errorf("bad permissions")
}
