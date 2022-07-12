package ships

import (
	"context"

	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Count(ctx context.Context) (int64, error) {
	return persistence.Count(ctx, "ships")
}

func CreateShip(ctx context.Context, ship entities.ShipEntity) error {
	return persistence.InsertOneToCollection(ctx, "ships", ship)
}

func GetShip(ctx context.Context, frn string) (entities.ShipEntity, error) {
	return persistence.FindOneByFrn[entities.ShipEntity](ctx, "ships", frn)
}

func ListShips(ctx context.Context, offset *int64, pageSize *int64) ([]entities.ShipEntity, error) {
	return persistence.List[entities.ShipEntity](ctx, "ships", options.Find().SetSkip(*offset).SetLimit(*pageSize))
}

func GetCargo(ctx context.Context, frn string) ([]entities.CargoEntity, error) {
	col, err := persistence.GetCollection("cargo")
	if err != nil {
		return nil, err
	}
	cur, err := col.Find(ctx, bson.M{"shipFrn": frn}, &options.FindOptions{})
	if err != nil {
		return nil, err
	}
	return persistence.DecodeCursor[entities.CargoEntity](ctx, cur)
}
