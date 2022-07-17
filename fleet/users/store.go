package users

import (
	"context"

	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Count(ctx context.Context) (int64, error) {
	return persistence.Count(ctx, "users")
}

func ListUsers(ctx context.Context, offset int64, pageSize int64, sort bson.D) ([]entities.UserEntity, error) {
	return persistence.List[entities.UserEntity](ctx, "users", options.Find().SetSkip(offset).SetLimit(pageSize).SetSort(sort))
}

func CreateUser(ctx context.Context, userEntity entities.UserEntity) error {
	return persistence.InsertOneToCollection(ctx, "users", userEntity)
}
