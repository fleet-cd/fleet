package cargo

import (
	"context"

	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CountCargo(ctx context.Context) (int64, error) {
	return persistence.Count(ctx, "cargo")
}

func ListCargo(ctx context.Context, offset *int64, pageSize *int64) ([]entities.CargoEntity, error) {
	return persistence.List[entities.CargoEntity](ctx, "cargo", options.Find().SetSkip(*offset).SetLimit(*pageSize))
}

func CreateCargo(ctx context.Context, cargo entities.CargoEntity) error {
	return persistence.InsertOneToCollection(ctx, "cargo", cargo)
}
