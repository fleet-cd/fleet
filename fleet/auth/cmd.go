package auth

import (
	"context"
	"time"

	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/frn"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
)

func CreateAdminUser(email string, password string) {
	hashedPW, err := GeneratePassword(password)
	if err != nil {
		panic(err)
	}
	now := time.Now()
	permission := common.NewPermissionBuilder().
		SetActions([]string{"*"}).
		SetFrn(common.PermissionFrn(frn.Generate("permission", "default").String())).
		SetCreatedAt(now).
		SetModifiedAt(now).
		SetName("admin access").
		SetNamespace("*").
		SetResourceType("*").
		Build()
	group := entities.NewGroupEntityBuilder().
		SetName("admins").
		SetCreatedAt(now).
		SetPermissions([]common.PermissionFrn{permission.Frn}).
		SetModifiedAt(now).Build()
	user := entities.NewUserEntityBuilder().
		SetFrn(common.UserFrn(frn.Generate("user", "default").String())).
		SetEmail(email).
		SetPassword(hashedPW).
		SetGroups([]string{group.Name}).
		SetCreatedAt(now).
		SetModifiedAt(now).Build()
	persistence.InsertOneToCollection(context.TODO(), "permissions", permission)
	persistence.InsertOneToCollection(context.TODO(), "groups", group)
	persistence.InsertOneToCollection(context.TODO(), "users", user)
}
