package server

import (
	"github.com/gin-gonic/gin"
	"github.com/tgs266/rest-gen/runtime/errors"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if ctx.Errors.Last() != nil {
			err := ctx.Errors.Last().Err
			if errors.IsStandardError(err) {
				err := err.(errors.StandardError)
				ctx.JSON(err.StatusCode, err)
			} else {
				ctx.JSON(500, err)
			}
		}
		ctx.Next()
	}
}
