package products

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/tgs266/fleet/fleet/panicker"
	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/fleet/utils"
	versionStore "github.com/tgs266/fleet/fleet/versions"
	"github.com/tgs266/fleet/frn"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/errors"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/products"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/versions"
	errs "github.com/tgs266/rest-gen/runtime/errors"
)

type ProductService struct{}

func (ss *ProductService) CreateProduct(body products.CreateProductRequest) (products.Product, error) {
	ns := utils.OrDefault(body.Namespace, "default")
	now := time.Now()
	prodFrn := frn.GenerateActual[common.ProductFrn]("ship", ns)

	product := products.NewProductBuilder().
		SetName(body.Name).
		SetNamespace(ns).
		SetCreatedAt(now).
		SetModifiedAt(now).
		SetEnvironment(body.Environment).
		SetFrn(prodFrn).Build()

	if err := CreateProduct(context.Background(), product); err != nil {
		return products.Product{}, err
	}

	return product, nil
}

func (ss *ProductService) ListProducts(offset *int64, pageSize *int64) (products.ListProductResponse, error) {
	total := panicker.AndPanic(Count(context.TODO())).GetOrPanicFunc(errs.NewInternalError)
	res := panicker.AndPanic(ListProducts(context.TODO(), offset, pageSize)).GetOrPanicFunc(errs.NewNotFound)
	results := utils.ConvertList(res, func(input products.Product) products.Product {
		return products.Product(input)
	})
	return products.NewListProductResponseBuilder().
		SetCount(len(results)).
		SetTotal(int(total)).
		SetItems(results).
		Build(), nil
}

func (ss *ProductService) GetProduct(expandVersions *bool, frn string) (products.GetProductResponse, error) {
	prodFrn := common.ProductFrn(frn)
	res := panicker.AndPanic(GetProduct(context.Background(), prodFrn)).GetOrPanicFunc(func(err error) error { return errors.NewProductNotFound(err, prodFrn) })
	vers := []versions.Version{}
	if *expandVersions {
		versionFrns := utils.Values(res.Versions)
		vers = panicker.AndPanic(versionStore.BatchGetVersions(context.TODO(), versionFrns)).GetOrPanicFunc(errs.NewInternalError)
	}
	response := products.GetProductResponse{
		Versions: vers,
		Product:  res,
	}
	return response, nil
}

func (ss *ProductService) AddVersion(body versions.CreateVersionRequest, pfrn string) (versions.Version, error) {
	prodFrn := common.ProductFrn(pfrn)

	product := panicker.AndPanic(GetProduct(context.TODO(), prodFrn)).GetOrPanicFunc(func(err error) error { return errors.NewProductNotFound(err, prodFrn) })
	now := time.Now()
	versionFrn := frn.GenerateActual[common.VersionFrn]("version", "default")

	version := versions.NewVersionBuilder().
		SetFrn(versionFrn).
		SetCreatedAt(now).
		SetModifiedAt(now).
		SetArtifactLocation(utils.OrDefault(body.ArtifactLocation, "")).
		SetVersion(body.Version).Build()

	panicker.CheckAndPanic(persistence.InsertOneToCollection(context.TODO(), "versions", version))
	if product.Versions == nil {
		product.Versions = map[string]common.VersionFrn{}
	}
	product.Versions[version.Version] = version.Frn
	panicker.CheckAndPanic(UpdateProduct(context.TODO(), prodFrn, product))
	return version, nil
}

func AddVersionArtifact(body *multipart.FileHeader, pfrn string, vfrn string) error {

	fileData := panicker.AndPanic(body.Open()).GetOrPanic()

	version := panicker.AndPanic(versionStore.GetVersion(context.TODO(), vfrn)).GetOrPanicFunc(errs.NewNotFound)
	product := panicker.AndPanic(GetProduct(context.TODO(), common.ProductFrn(pfrn))).GetOrPanicFunc(errs.NewNotFound)

	folderPath := filepath.Join("artifacts", string(product.Frn), string(version.Frn))
	if runtime.GOOS == "windows" {
		folderPath = strings.ReplaceAll(folderPath, ":", ".")
	}
	filePath := filepath.Join(folderPath, "artifact"+filepath.Ext(body.Filename))

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		panicker.CheckAndPanic(os.MkdirAll(folderPath, os.ModePerm))
		panicker.AndPanic(os.Create(filePath)).GetOrPanic()
	}

	out := panicker.AndPanic(os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)).GetOrPanic()
	defer out.Close()
	panicker.AndPanic(io.Copy(out, fileData)).GetOrPanic()

	version.ArtifactLocation = filePath
	version.ModifiedAt = time.Now()
	panicker.CheckAndPanic(versionStore.UpdateVersion(context.TODO(), version.Frn, version))

	return nil
}
