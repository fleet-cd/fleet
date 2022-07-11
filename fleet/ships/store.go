package ships

import (
	"context"

	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateShip(ctx context.Context, ship entities.ShipEntity) error {
	return persistence.InsertOneToCollection(ctx, "ships", ship)
}

func GetShip(ctx context.Context, frn string) (entities.ShipEntity, error) {
	return persistence.FindOneByFrn[entities.ShipEntity](ctx, "ships", frn)

}

func ListShips(ctx context.Context) ([]entities.ShipEntity, error) {
	col := persistence.GetCollection("ships")
	cur, err := col.Find(ctx, bson.M{}, &options.FindOptions{})
	if err != nil {
		return nil, err
	}
	return persistence.DecodeCursor[entities.ShipEntity](ctx, cur)
}
