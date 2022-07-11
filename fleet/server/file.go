package server

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	errs "github.com/tgs266/fleet/rest-gen/generated/com/fleet/errors"
	"github.com/tgs266/rest-gen/runtime/errors"
)

func HandleManifestUpload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")

	if err != nil {
		ctx.AbortWithError(400, errors.NewInvalidArgumentError(err))
		return
	}
	shipFrn := ctx.Param("shipFrn")

	extension := filepath.Ext(file.Filename)
	if extension != ".yaml" && extension != ".yml" {
		ctx.AbortWithError(400, errs.NewInvalidFileType(err, []string{".yaml", ".yml"}))
		return
	}
	newFileName := shipFrn + ".yaml"

	// The file is received, so let's save it
	if err := ctx.SaveUploadedFile(file, "manifests/"+newFileName); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.StandardError{
			Name:       "Internal",
			Code:       "INTERNAL",
			StatusCode: 500,
			ErrorId:    uuid.NewString(),
		})
		return
	}

	ctx.Status(200)
}
