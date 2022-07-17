package securedauth

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	fleetAuth "github.com/tgs266/fleet/fleet/auth"
	"github.com/tgs266/fleet/fleet/panicker"
	"github.com/tgs266/fleet/fleet/utils"
	"github.com/tgs266/fleet/frn"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/auth"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
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
	groupEntity := entities.NewGroupEntityBuilder().
		SetName(body.Name).
		SetPermissions(body.Permissions).
		SetCreatedAt(now).
		SetModifiedAt(now).Build()

	panicker.CheckAndPanic(CreateGroup(ctx, groupEntity))
	return groupEntity, nil
}

func (service *SecuredAuthService) GetGroup(ctx *gin.Context, expandPermissions bool, name string, token authentication.Token) (auth.GetGroupResponse, error) {
	panicker.AndPanic(fleetAuth.CanINoNamepsace(token, "group", fleetAuth.ACTION_VIEW)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	group := panicker.AndPanic(GetGroup(ctx, name)).GetOrPanicFunc(errors.NewNotFound)
	perms := []common.Permission{}
	if expandPermissions {
		perms = panicker.AndPanic(BactchGetPermissions(ctx, group.Permissions)).GetOrPanicFunc(errors.NewNotFound)
	}
	return auth.GetGroupResponse{
		Group:               group,
		ExpandedPermissions: perms,
	}, nil
}

func (service *SecuredAuthService) DeleteGroupPermission(ctx *gin.Context, name string, permissionName string, token authentication.Token) error {
	panicker.AndPanic(fleetAuth.CanINoNamepsace(token, "group", fleetAuth.ACTION_EDIT)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	group := panicker.AndPanic(GetGroup(ctx, name)).GetOrPanicFunc(errors.NewNotFound)
	fmt.Println(group.Permissions)
	fmt.Println(utils.FindAndRemove(group.Permissions, common.PermissionFrn(permissionName)))
	group.Permissions = utils.FindAndRemove(group.Permissions, common.PermissionFrn(permissionName))
	fmt.Println(group.Permissions)
	panicker.CheckAndPanic(UpdateGroup(ctx, name, group))
	return nil
}

func (service *SecuredAuthService) ListPermissions(ctx *gin.Context, sort string, token authentication.Token) ([]common.Permission, error) {
	panicker.AndPanic(fleetAuth.CanI(token, "permission", "*", fleetAuth.ACTION_VIEW)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	perms, err := ListPermissions(ctx, utils.GetSortMap(sort))
	if err != nil {
		return nil, errors.NewNotFound(err)
	}
	return perms, nil
}

func (service *SecuredAuthService) CreatePermission(ctx *gin.Context, body auth.CreatePermissionRequest, token authentication.Token) (common.Permission, error) {
	panicker.AndPanic(fleetAuth.CanI(token, "permission", "*", fleetAuth.ACTION_CREATE)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	now := time.Now()
	entity := common.NewPermissionBuilder().
		SetActions(body.Actions).
		SetCreatedAt(now).
		SetModifiedAt(now).
		SetFrn(frn.GenerateActual[common.PermissionFrn]("permission", "default")).
		SetName(body.Name).
		SetNamespace(body.Namespace).
		SetResourceType(body.ResourceType).Build()
	panicker.CheckAndPanic(CreatePermission(ctx, entity))
	return entity, nil
}
