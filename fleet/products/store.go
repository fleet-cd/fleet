package products

import (
	"context"

	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Count(ctx context.Context) (int64, error) {
	return persistence.Count(ctx, "products")
}

func CreateProduct(ctx context.Context, product entities.ProductEntity) error {
	return persistence.InsertOneToCollection(ctx, "products", product)
}

func GetProduct(ctx context.Context, frn common.ProductFrn) (entities.ProductEntity, error) {
	return persistence.FindOneByFrn[entities.ProductEntity](ctx, "products", string(frn))
}

func ListProducts(ctx context.Context, offset *int64, pageSize *int64) ([]entities.ProductEntity, error) {
	return persistence.List[entities.ProductEntity](ctx, "products", options.Find().SetSkip(*offset).SetLimit(*pageSize))
}

func GetCargo(ctx context.Context, frn common.ProductFrn) ([]entities.CargoEntity, error) {
	col, err := persistence.GetCollection("cargo")
	if err != nil {
		return nil, err
	}
	cur, err := col.Find(ctx, bson.M{"productFrn": frn}, &options.FindOptions{})
	if err != nil {
		return nil, err
	}
	return persistence.DecodeCursor[entities.CargoEntity](ctx, cur)
}
