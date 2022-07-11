package cargo

import (
	"context"
	"time"

	productService "github.com/tgs266/fleet/fleet/products"
	"github.com/tgs266/fleet/fleet/ships"
	"github.com/tgs266/fleet/frn"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/cargo"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/errors"
)

type CargoService struct{}

func (service *CargoService) CreateCargo(body cargo.CreateCargoRequest) (cargo.Cargo, error) {
	now := time.Now()
	frn := frn.Generate("cargo", "main").String()

	if _, err := productService.GetProduct(context.TODO(), body.ProductFrn); err != nil {
		return cargo.Cargo{}, errors.NewProductNotFound(err, body.ProductFrn)
	}

	if _, err := ships.GetShip(context.TODO(), body.ShipFrn); err != nil {
		return cargo.Cargo{}, errors.NewShipNotFound(err, body.ShipFrn)
	}

	entity := entities.NewCargoEntityBuilder().
		SetFrn(frn).
		SetShipFrn(body.ShipFrn).
		SetProductFrn(body.ProductFrn).
		SetCreatedAt(now).
		SetModifiedAt(now)

	if err := CreateCargo(context.Background(), entity.Build()); err != nil {
		return cargo.Cargo{}, err
	}

	cargo := cargo.NewCargoBuilder().
		SetFrn(entity.Frn).
		SetProductFrn(entity.ProductFrn).
		SetShipFrn(entity.ShipFrn).
		SetCreatedAt(entity.CreatedAt).
		SetModifiedAt(entity.ModifiedAt)
	return cargo.Build(), nil
}
