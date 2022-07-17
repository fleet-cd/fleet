package products

import (
	"context"
	"time"

	"github.com/tgs266/fleet/fleet/utils"
	"github.com/tgs266/fleet/frn"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/cargo"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/errors"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/products"
)

type ProductService struct{}

func (ss *ProductService) CreateProduct(body products.CreateProductRequest) (products.Product, error) {
	ns := utils.OrDefault(body.Namespace, "default")
	now := time.Now()
	prodFrn := frn.GenerateActual[common.ProductFrn]("ship", ns)

	entity := entities.NewProductEntityBuilder().
		SetName(body.Name).
		SetNamespace(ns).
		SetCreatedAt(now).
		SetModifiedAt(now).
		SetFrn(prodFrn)

	if err := CreateProduct(context.Background(), entity.Build()); err != nil {
		return products.Product{}, err
	}

	product := products.Product(entity)
	return product, nil
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

	results := utils.ConvertList(res, func(input entities.ProductEntity) products.Product {
		return products.Product(input)
	})

	return products.NewListProductResponseBuilder().
		SetCount(len(results)).
		SetTotal(int(total)).
		SetItems(results).
		Build(), nil
}

func (ss *ProductService) GetProduct(frn string) (products.Product, error) {
	prodFrn := common.ProductFrn(frn)
	res, err := GetProduct(context.Background(), prodFrn)
	if err != nil {
		return products.Product{}, errors.NewProductNotFound(err, prodFrn)
	}

	result := products.Product(res)
	return result, nil
}

func (ss *ProductService) GetCargo(frn string) ([]cargo.Cargo, error) {
	prodFrn := common.ProductFrn(frn)
	res, err := GetCargo(context.Background(), prodFrn)
	if err != nil {
		return nil, errors.NewProductNotFound(err, prodFrn)
	}

	results := utils.ConvertList(res, func(input entities.CargoEntity) cargo.Cargo {
		return cargo.Cargo(input)
	})

	return results, nil
}
