package cargo

import (
	"context"

	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
)

func CreateCargo(ctx context.Context, cargo entities.CargoEntity) error {
	return persistence.InsertOneToCollection(ctx, "cargo", cargo)
}
