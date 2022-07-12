package ships

import (
	"context"
	"time"

	"github.com/tgs266/fleet/fleet/utils"
	"github.com/tgs266/fleet/frn"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/cargo"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/errors"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/ships"
)

type ShipService struct{}

func (ss *ShipService) CreateShip(body ships.CreateShipRequest) (ships.Ship, error) {
	ns := utils.OrDefault(body.Namespace, "default")
	now := time.Now()
	shipFrn := frn.Generate("ship", ns).String()

	entity := entities.NewShipEntityBuilder().
		SetName(body.Name).
		SetNamespace(ns).
		SetCreatedAt(now).
		SetModifiedAt(now).
		SetFrn(shipFrn)

	if err := CreateShip(context.Background(), entity.Build()); err != nil {
		return ships.Ship{}, err
	}

	ship := ships.NewShipBuilder().
		SetFrn(entity.Frn).
		SetName(entity.Name).
		SetNamespace(entity.Namespace).
		SetCreatedAt(entity.CreatedAt).
		SetModifiedAt(entity.ModifiedAt)
	return ship.Build(), nil
}

func (ss *ShipService) ListShips(offset *int64, pageSize *int64) (ships.ListShipsResponse, error) {
	total, err := Count(context.TODO())
	if err != nil {
		return ships.ListShipsResponse{}, err
	}
	res, err := ListShips(context.TODO(), offset, pageSize)
	if err != nil {
		return ships.ListShipsResponse{}, err
	}

	results := utils.ConvertList(res, func(input entities.ShipEntity) ships.Ship {
		return ships.NewShipBuilder().
			SetFrn(input.Frn).
			SetName(input.Name).
			SetNamespace(input.Namespace).
			SetCreatedAt(input.CreatedAt).
			SetModifiedAt(input.ModifiedAt).
			Build()
	})

	return ships.NewListShipsResponseBuilder().
		SetTotal(int(total)).
		SetCount(len(results)).
		SetItems(results).
		Build(), nil
}

func (ss *ShipService) GetShip(shipFrn string) (ships.Ship, error) {

	res, err := GetShip(context.Background(), shipFrn)
	if err != nil {
		return ships.Ship{}, errors.NewShipNotFound(err, shipFrn)
	}

	shp := ships.NewShipBuilder().
		SetFrn(shipFrn).
		SetName(res.Name).
		SetNamespace(res.Namespace).
		SetCreatedAt(res.CreatedAt).
		SetModifiedAt(res.ModifiedAt).Build()

	return shp, nil
}

func (ss *ShipService) GetCargo(frn string) ([]cargo.Cargo, error) {

	res, err := GetCargo(context.Background(), frn)
	if err != nil {
		return nil, errors.NewShipNotFound(err, frn)
	}

	results := []cargo.Cargo{}
	for _, r := range res {
		results = append(results, cargo.NewCargoBuilder().
			SetFrn(r.Frn).
			SetProductFrn(r.ProductFrn).
			SetShipFrn(r.ShipFrn).
			SetCreatedAt(r.CreatedAt).
			SetModifiedAt(r.ModifiedAt).
			Build(),
		)
	}

	return results, nil
}
