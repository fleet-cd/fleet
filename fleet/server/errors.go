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
			er := errors.GetError(err)
			ctx.JSON(er.Code(), err)
		}
	}
}

func RecoveryMiddleware(handler func(ctx *gin.Context, err error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				handler(ctx, err.(error))
			}
		}()
		ctx.Next()
	}
}
