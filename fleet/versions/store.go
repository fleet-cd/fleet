package version

import (
	"context"

	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/versions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateVersion(ctx context.Context, version versions.Version) error {
	return persistence.InsertOneToCollection(ctx, "versions", version)
}

func BatchGetVersions(ctx context.Context, vers []common.VersionFrn) ([]versions.Version, error) {
	return persistence.ListQuery[versions.Version](ctx, "versions", bson.M{"frn": bson.M{"$in": vers}}, options.Find())
}

func GetVersion(ctx context.Context, vers string) (versions.Version, error) {
	return persistence.FindOneByFrn[versions.Version](ctx, "versions", string(vers))
}

func UpdateVersion(ctx context.Context, vers common.VersionFrn, version versions.Version) error {
	return persistence.UpdateOne(ctx, "versions", bson.M{"frn": vers}, version)
}
