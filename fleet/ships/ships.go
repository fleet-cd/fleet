package ships

import (
	"context"
	"time"

	"github.com/tgs266/fleet/fleet/auth"
	"github.com/tgs266/fleet/fleet/panicker"
	"github.com/tgs266/fleet/fleet/utils"
	"github.com/tgs266/fleet/fleet/validation"
	"github.com/tgs266/fleet/frn"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/cargo"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/errors"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/ships"
	"github.com/tgs266/rest-gen/runtime/authentication"
)

type ShipService struct {
}

func (ss *ShipService) CreateShip(body ships.CreateShipRequest, token authentication.Token) (ships.Ship, error) {
	ns := utils.OrDefault(body.Namespace, "default")
	panicker.AndPanic(auth.CanI(token, "ship", ns, auth.ACTION_CREATE)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	panicker.OnFalse(validation.ValidateName(ns), errors.NewInvalidName(nil, ns))
	panicker.OnFalse(validation.ValidateName(body.Name), errors.NewInvalidName(nil, body.Name))
	now := time.Now()
	shipFrn := frn.GenerateActual[common.ShipFrn]("ship", ns)

	entity := entities.NewShipEntityBuilder().
		SetName(body.Name).
		SetNamespace(ns).
		SetCreatedAt(now).
		SetModifiedAt(now).
		SetFrn(shipFrn).
		SetTags(body.Tags)

	if err := CreateShip(context.Background(), entity.Build()); err != nil {
		return ships.Ship{}, err
	}

	ship := ships.Ship(entity)
	return ship, nil
}

func (ss *ShipService) ListShips(offset *int64, pageSize *int64, sort string, token authentication.Token) (ships.ListShipsResponse, error) {

	request := panicker.AndPanic(auth.WhatCanIList(token, "ship")).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	total, err := Count(context.TODO(), request)
	if err != nil {
		return ships.ListShipsResponse{}, err
	}
	sortMap := utils.GetSortMap(sort)

	res, err := ListShips(context.Background(), request, utils.OrDefault(offset, 0), utils.OrDefault(pageSize, 0), sortMap)
	if err != nil {
		return ships.ListShipsResponse{}, err
	}

	results := utils.ConvertList(res, func(input entities.ShipEntity) ships.Ship {
		return ships.Ship(input)
	})

	return ships.NewListShipsResponseBuilder().
		SetTotal(int(total)).
		SetCount(len(results)).
		SetItems(results).
		Build(), nil
}

func (ss *ShipService) GetShip(shipFrn string, token authentication.Token) (ships.Ship, error) {
	panicker.AndPanic(auth.CanIFrn(token, shipFrn, auth.ACTION_VIEW)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	frn := common.ShipFrn(shipFrn)
	res, err := GetShip(context.Background(), frn)
	if err != nil {
		return ships.Ship{}, errors.NewShipNotFound(err, frn)
	}

	shp := ships.Ship(res)

	return shp, nil
}

func (ss *ShipService) DeleteShip(shipFrn string, token authentication.Token) (common.ShipFrn, error) {
	panicker.AndPanic(auth.CanIFrn(token, shipFrn, auth.ACTION_DELETE)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	frn := common.ShipFrn(shipFrn)
	err := DeleteShip(context.Background(), frn)
	if err != nil {
		return "", errors.NewShipNotFound(err, frn)
	}

	return frn, nil
}

func (ss *ShipService) GetCargo(frn string, token authentication.Token) ([]cargo.Cargo, error) {
	panicker.AndPanic(auth.CanI(token, "cargo", "*", auth.ACTION_VIEW)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	shipFrn := common.ShipFrn(frn)
	res, err := GetCargo(context.Background(), shipFrn)
	if err != nil {
		return nil, errors.NewShipNotFound(err, shipFrn)
	}

	results := utils.ConvertList(res, func(input entities.CargoEntity) cargo.Cargo {
		return cargo.Cargo(input)
	})

	return results, nil
}
