package cargo

type CargoService struct{}

// func (service *CargoService) CreateCargo(body cargo.CreateCargoRequest) (cargo.Cargo, error) {
// 	now := time.Now()
// 	frn := common.CargoFrn(frn.Generate("cargo", "default").String())

// 	if _, err := productService.GetProduct(context.TODO(), body.ProductFrn); err != nil {
// 		return cargo.Cargo{}, errors.NewProductNotFound(err, body.ProductFrn)
// 	}

// 	if _, err := ships.GetShip(context.TODO(), body.ShipFrn, bson.M{}); err != nil {
// 		return cargo.Cargo{}, errors.NewShipNotFound(err, body.ShipFrn)
// 	}

// 	entity := entities.NewCargoEntityBuilder().
// 		SetFrn(frn).
// 		SetProductFrn(body.ProductFrn).
// 		SetCreatedAt(now).
// 		SetModifiedAt(now)

// 	if err := CreateCargo(context.Background(), entity.Build()); err != nil {
// 		return cargo.Cargo{}, err
// 	}

// 	cargo := cargo.NewCargoBuilder().
// 		SetFrn(entity.Frn).
// 		SetProductFrn(entity.ProductFrn).
// 		SetCreatedAt(entity.CreatedAt).
// 		SetModifiedAt(entity.ModifiedAt)
// 	return cargo.Build(), nil
// }

// func (service *CargoService) ListCargo(offset *int64, pageSize *int64) (cargo.ListCargoResponse, error) {
// 	total, err := CountCargo(context.TODO())
// 	if err != nil {
// 		return cargo.ListCargoResponse{}, err
// 	}

// 	items, err := ListCargo(context.TODO(), offset, pageSize)
// 	if err != nil {
// 		return cargo.ListCargoResponse{}, err
// 	}
// 	results := []cargo.Cargo{}
// 	for _, r := range items {
// 		results = append(results, cargo.NewCargoBuilder().
// 			SetFrn(r.Frn).
// 			SetShipFrn(r.ShipFrn).
// 			SetProductFrn(r.ProductFrn).
// 			SetCreatedAt(r.CreatedAt).
// 			SetModifiedAt(r.ModifiedAt).
// 			Build(),
// 		)
// 	}

// 	return cargo.ListCargoResponse{
// 		Total: int(total),
// 		Count: len(results),
// 		Items: results,
// 	}, nil
// }
