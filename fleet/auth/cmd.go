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
	permission := entities.NewPermissionEntityBuilder().
		SetActions([]string{"*"}).
		SetCreatedAt(now).
		SetModifiedAt(now).
		SetNamespace("*").
		SetResourceType("*").
		Build()
	group := entities.NewGroupEntityBuilder().
		SetName("admins").
		SetCreatedAt(now).
		SetPermissions([]entities.PermissionEntity{permission}).
		SetModifiedAt(now).Build()
	user := entities.NewUserEntityBuilder().
		SetFrn(common.UserFrn(frn.Generate("user", "default").String())).
		SetEmail(email).
		SetPassword(hashedPW).
		SetGroups([]string{group.Name}).
		SetCreatedAt(now).
		SetModifiedAt(now).Build()
	persistence.InsertOneToCollection(context.TODO(), "groups", group)
	persistence.InsertOneToCollection(context.TODO(), "users", user)
}
