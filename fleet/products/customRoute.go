package products

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tgs266/rest-gen/runtime/errors"
)

func HandleAddVersionArtifact(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.AbortWithError(500, errors.NewInvalidArgumentError(err))
	}
	productFrn := ctx.Param("productFrn")
	versionFrn := ctx.Param("versionFrn")
	if err = AddVersionArtifact(file, productFrn, versionFrn); err != nil {
		ctx.AbortWithError(500, err)
		return
	} else {
		ctx.Status(http.StatusOK)
		return
	}
}
