package products

import (
	"context"

	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/products"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Count(ctx context.Context) (int64, error) {
	return persistence.Count(ctx, "products")
}

func CreateProduct(ctx context.Context, product products.Product) error {
	return persistence.InsertOneToCollection(ctx, "products", product)
}

func GetProduct(ctx context.Context, frn common.ProductFrn) (products.Product, error) {
	return persistence.FindOneByFrn[products.Product](ctx, "products", string(frn))
}

func PutProduct(ctx context.Context, frn common.ProductFrn, product products.Product) error {
	return persistence.UpdateOne[products.Product](ctx, "products", bson.M{"frn": frn}, product)
}

func ListProducts(ctx context.Context, offset *int64, pageSize *int64) ([]products.Product, error) {
	return persistence.List[products.Product](ctx, "products", options.Find().SetSkip(*offset).SetLimit(*pageSize))
}

func UpdateProduct(ctx context.Context, frn common.ProductFrn, product products.Product) error {
	return persistence.UpdateOne(ctx, "products", bson.M{"frn": frn}, product)
}

// func GetCargo(ctx context.Context, frn common.ProductFrn) ([]entities.CargoEntity, error) {
// 	col, err := persistence.GetCollection("cargo")
// 	if err != nil {
// 		return nil, err
// 	}
// 	cur, err := col.Find(ctx, bson.M{"productFrn": frn}, &options.FindOptions{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return persistence.DecodeCursor[entities.CargoEntity](ctx, cur)
// }
