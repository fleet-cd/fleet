package products

import (
	"context"

	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateProduct(ctx context.Context, product entities.ProductEntity) error {
	return persistence.InsertOneToCollection(ctx, "products", product)
}

func GetProduct(ctx context.Context, frn string) (entities.ProductEntity, error) {
	return persistence.FindOneByFrn[entities.ProductEntity](ctx, "products", frn)
}

func ListProducts(ctx context.Context) ([]entities.ProductEntity, error) {
	col := persistence.GetCollection("products")
	cur, err := col.Find(ctx, bson.M{}, &options.FindOptions{})
	if err != nil {
		return nil, err
	}
	return persistence.DecodeCursor[entities.ProductEntity](ctx, cur)
}

func GetCargo(ctx context.Context, frn string) ([]entities.CargoEntity, error) {
	col := persistence.GetCollection("cargo")
	cur, err := col.Find(ctx, bson.M{"productFrn": frn}, &options.FindOptions{})
	if err != nil {
		return nil, err
	}
	return persistence.DecodeCursor[entities.CargoEntity](ctx, cur)
}
