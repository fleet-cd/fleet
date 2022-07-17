package utils

import "github.com/gin-gonic/gin"

func GetGinRecover(ctx *gin.Context) func() {
	return func() {
		if r := recover(); r != nil {
			ctx.AbortWithError(500, r.(error))
		}
	}
}
