package auth

import (
	"context"

	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUserByFrn(userFrn common.UserFrn) (entities.UserEntity, error) {
	return persistence.FindOneByFrn[entities.UserEntity](context.TODO(), "users", string(userFrn))
}

func GetUserByEmail(email string) (entities.UserEntity, error) {
	return persistence.FindOne[entities.UserEntity](context.TODO(), "users", bson.M{"email": email})
}

func ListGroups(ctx context.Context) ([]entities.GroupEntity, error) {
	return persistence.List[entities.GroupEntity](ctx, "groups", &options.FindOptions{})
}
