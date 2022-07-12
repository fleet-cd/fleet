package products

import (
	"context"
	"time"

	"github.com/tgs266/fleet/fleet/utils"
	"github.com/tgs266/fleet/frn"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/cargo"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/errors"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/products"
)

type ProductService struct{}

func (ss *ProductService) CreateProduct(body products.CreateProductRequest) (products.Product, error) {
	ns := utils.OrDefault(body.Namespace, "default")
	now := time.Now()
	shipFrn := frn.Generate("product", ns).String()

	entity := entities.NewProductEntityBuilder().
		SetName(body.Name).
		SetNamespace(ns).
		SetCreatedAt(now).
		SetModifiedAt(now).
		SetFrn(shipFrn)

	if err := CreateProduct(context.Background(), entity.Build()); err != nil {
		return products.Product{}, err
	}

	product := products.NewProductBuilder().
		SetFrn(entity.Frn).
		SetName(entity.Name).
		SetNamespace(entity.Namespace).
		SetCreatedAt(entity.CreatedAt).
		SetModifiedAt(entity.ModifiedAt)
	return product.Build(), nil
}

func (ss *ProductService) ListProducts(offset *int64, pageSize *int64) (products.ListProductResponse, error) {
	total, err := Count(context.TODO())
	if err != nil {
		return products.ListProductResponse{}, err
	}
	res, err := ListProducts(context.TODO(), offset, pageSize)
	if err != nil {
		return products.ListProductResponse{}, err
	}

	results := utils.ConvertList[entities.ProductEntity, products.Product](res, func(input entities.ProductEntity) products.Product {
		return products.NewProductBuilder().
			SetFrn(input.Frn).
			SetName(input.Name).
			SetNamespace(input.Namespace).
			SetCreatedAt(input.CreatedAt).
			SetModifiedAt(input.ModifiedAt).
			Build()
	})

	return products.NewListProductResponseBuilder().
		SetCount(len(results)).
		SetTotal(int(total)).
		SetItems(results).
		Build(), nil
}

func (ss *ProductService) GetProduct(frn string) (products.Product, error) {

	res, err := GetProduct(context.Background(), frn)
	if err != nil {
		return products.Product{}, errors.NewProductNotFound(err, frn)
	}

	result := products.NewProductBuilder().
		SetFrn(frn).
		SetName(res.Name).
		SetNamespace(res.Namespace).
		SetCreatedAt(res.CreatedAt).
		SetModifiedAt(res.ModifiedAt).Build()

	return result, nil
}

func (ss *ProductService) GetCargo(frn string) ([]cargo.Cargo, error) {

	res, err := GetCargo(context.Background(), frn)
	if err != nil {
		return nil, errors.NewProductNotFound(err, frn)
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
