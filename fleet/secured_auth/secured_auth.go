package securedauth

import (
	"time"

	"github.com/gin-gonic/gin"
	fleetAuth "github.com/tgs266/fleet/fleet/auth"
	"github.com/tgs266/fleet/fleet/panicker"
	"github.com/tgs266/fleet/fleet/utils"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/auth"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"github.com/tgs266/rest-gen/runtime/authentication"
	"github.com/tgs266/rest-gen/runtime/errors"
)

type SecuredAuthService struct{}

func (service *SecuredAuthService) ListGroups(ctx *gin.Context, sort string, token authentication.Token) ([]entities.GroupEntity, error) {
	panicker.AndPanic(fleetAuth.CanI(token, "group", "*", fleetAuth.ACTION_VIEW)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	groups, err := ListGroups(ctx, utils.GetSortMap(sort))
	if err != nil {
		return nil, errors.NewNotFound(err)
	}
	return groups, nil
}

func (service *SecuredAuthService) CreateGroup(ctx *gin.Context, body auth.CreateGroupRequest, token authentication.Token) (entities.GroupEntity, error) {
	panicker.AndPanic(fleetAuth.CanI(token, "group", "*", fleetAuth.ACTION_CREATE)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	now := time.Now()

	permissions := utils.ConvertList(body.Permissions, func(input auth.CreatePermissionRequest) entities.PermissionEntity {
		return entities.NewPermissionEntityBuilder().
			SetActions(input.Actions).
			SetNamespace(input.Namespace).
			SetResourceType(input.ResourceType).
			SetCreatedAt(now).
			SetModifiedAt(now).
			Build()
	})

	groupEntity := entities.NewGroupEntityBuilder().
		SetName(body.Name).
		SetPermissions(permissions).
		SetCreatedAt(now).
		SetModifiedAt(now).Build()

	panicker.CheckAndPanic(CreateGroup(ctx, groupEntity))
	return groupEntity, nil
}

func (service *SecuredAuthService) GetGroup(ctx *gin.Context, name string, token authentication.Token) (entities.GroupEntity, error) {
	panicker.AndPanic(fleetAuth.CanINoNamepsace(token, "group", fleetAuth.ACTION_VIEW)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	group := panicker.AndPanic(GetGroup(ctx, name)).GetOrPanicFunc(errors.NewNotFound)
	return group, nil
}

func (service *SecuredAuthService) DeleteGroupPermission(ctx *gin.Context, name string, permissionIdx int, token authentication.Token) error {
	panicker.AndPanic(fleetAuth.CanINoNamepsace(token, "group", fleetAuth.ACTION_EDIT)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	group := panicker.AndPanic(GetGroup(ctx, name)).GetOrPanicFunc(errors.NewNotFound)
	group.Permissions = utils.RemoveByIdx(group.Permissions, permissionIdx)
	panicker.CheckAndPanic(UpdateGroup(ctx, name, group))
	return nil
}

func (service *SecuredAuthService) AddGroupPermission(ctx *gin.Context, body auth.CreatePermissionRequest, groupName string, token authentication.Token) error {
	panicker.AndPanic(fleetAuth.CanI(token, "permission", "*", fleetAuth.ACTION_CREATE)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	now := time.Now()
	entity := entities.NewPermissionEntityBuilder().
		SetActions(body.Actions).
		SetCreatedAt(now).
		SetModifiedAt(now).
		SetNamespace(body.Namespace).
		SetResourceType(body.ResourceType).Build()

	group := panicker.AndPanic(GetGroup(ctx, groupName)).GetOrPanicFunc(errors.NewNotFound)
	group.Permissions = append(group.Permissions, entity)
	panicker.CheckAndPanic(UpdateGroup(ctx, groupName, group))
	return nil
}
