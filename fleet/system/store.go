package system

import (
	"context"

	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/system"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListEnvs(ctx context.Context, sort bson.D) ([]system.Environment, error) {
	return persistence.ListQuery[system.Environment](ctx, "environments", bson.M{}, options.Find().SetSort(sort))
}

func CreateEnv(ctx context.Context, env system.Environment) error {
	return persistence.InsertOneToCollection(ctx, "environments", env)
}
