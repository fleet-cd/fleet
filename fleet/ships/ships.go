package ships

import (
	"context"
	"time"

	"github.com/tgs266/fleet/fleet/utils"
	"github.com/tgs266/fleet/frn"
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

func (ss *ShipService) ListShips() ([]ships.Ship, error) {

	res, err := ListShips(context.Background())
	if err != nil {
		return nil, err
	}

	resultShips := []ships.Ship{}
	for _, r := range res {
		resultShips = append(resultShips, ships.NewShipBuilder().
			SetFrn(r.Frn).
			SetName(r.Name).
			SetNamespace(r.Namespace).
			SetCreatedAt(r.CreatedAt).
			SetModifiedAt(r.ModifiedAt).
			Build(),
		)
	}

	return resultShips, nil
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
