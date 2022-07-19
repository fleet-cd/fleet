package ships

import (
	"context"
	"fmt"
	"time"

	"github.com/tgs266/fleet/fleet/auth"
	"github.com/tgs266/fleet/fleet/git"
	"github.com/tgs266/fleet/fleet/kubernetes"
	"github.com/tgs266/fleet/fleet/panicker"
	"github.com/tgs266/fleet/fleet/utils"
	"github.com/tgs266/fleet/fleet/validation"
	"github.com/tgs266/fleet/frn"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/errors"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/ships"
	"github.com/tgs266/rest-gen/runtime/authentication"
	errs "github.com/tgs266/rest-gen/runtime/errors"
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

	entity := ships.NewShipBuilder().
		SetName(body.Name).
		SetNamespace(ns).
		SetCreatedAt(now).
		SetModifiedAt(now).
		SetFrn(shipFrn).
		SetTags(body.Tags).SetSource(body.Source).Build()

	if err := kubernetes.InitDeployment(entity); err != nil {
		fmt.Println(err)
		panic(err)
	}

	if err := CreateShip(context.Background(), entity); err != nil {
		return ships.Ship{}, err
	}

	return entity, nil
}

func (ss *ShipService) ListShips(offset *int64, pageSize *int64, sort string, token authentication.Token) (ships.ListShipsResponse, error) {
	request := panicker.AndPanic(auth.WhatCanIView(token, "ship")).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	total := panicker.AndPanic(Count(context.TODO(), request)).GetOr(0)
	sortMap := utils.GetSortMap(sort)
	res := panicker.AndPanic(
		ListShips(context.Background(), request, utils.OrDefault(offset, 0), utils.OrDefault(pageSize, 0), sortMap),
	).GetOrPanicFunc(errs.NewNotFound)
	return ships.NewListShipsResponseBuilder().
		SetTotal(int(total)).
		SetCount(len(res)).
		SetItems(res).
		Build(), nil
}

func (ss *ShipService) GetShip(shipFrn string, token authentication.Token) (ships.Ship, error) {
	constraint := panicker.AndPanic(auth.WhatCanIView(token, "ship")).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	frn := common.ShipFrn(shipFrn)
	res, err := GetShip(context.Background(), frn, constraint)
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

func (ss *ShipService) GetConfig(frn string, token authentication.Token) (ships.ConfigResponse, error) {
	constraint := panicker.AndPanic(auth.WhatCanIView(token, "ship")).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	shipFrn := common.ShipFrn(frn)
	res, err := GetShip(context.Background(), shipFrn, constraint)
	if err != nil {
		return ships.ConfigResponse{}, errors.NewShipNotFound(err, shipFrn)
	}

	cfg := panicker.AndPanic(git.GetConfig(res.Source)).GetOrPanicFunc(errs.NewInternalError)

	return ships.ConfigResponse{
		Body: cfg,
	}, nil
}
