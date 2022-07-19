package ships

import (
	"context"

	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/products"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/ships"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Count(ctx context.Context, query bson.M) (int64, error) {
	return persistence.CountQuery(ctx, "ships", query)
}

func CreateShip(ctx context.Context, ship ships.Ship) error {
	return persistence.InsertOneToCollection(ctx, "ships", ship)
}

func GetShip(ctx context.Context, frn common.ShipFrn, constraint bson.M) (ships.Ship, error) {
	constraint["frn"] = string(frn)
	return persistence.FindOne[ships.Ship](ctx, "ships", constraint)
}

func DeleteShip(ctx context.Context, frn common.ShipFrn) error {
	return persistence.DeleteOneByFrn[ships.Ship](ctx, "ships", string(frn))
}

func ListShips(ctx context.Context, request bson.M, offset int64, pageSize int64, sort bson.D) ([]ships.Ship, error) {
	return persistence.ListQuery[ships.Ship](ctx, "ships", request, options.Find().SetSkip(offset).SetLimit(pageSize).SetSort(sort))
}

// func GetCargo(ctx context.Context, frn common.ShipFrn) ([]entities.CargoEntity, error) {
// 	col, err := persistence.GetCollection("cargo")
// 	if err != nil {
// 		return nil, err
// 	}
// 	cur, err := col.Find(ctx, bson.M{"shipFrn": frn}, &options.FindOptions{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return persistence.DecodeCursor[entities.CargoEntity](ctx, cur)
// }

func BatchGetProducts(ctx context.Context, prods []common.ProductFrn) ([]products.Product, error) {
	col, err := persistence.GetCollection("cargo")
	if err != nil {
		return nil, err
	}
	cur, err := col.Find(ctx, bson.M{"productFrn": bson.M{"$in": prods}}, &options.FindOptions{})
	if err != nil {
		return nil, err
	}
	return persistence.DecodeCursor[products.Product](ctx, cur)
}
