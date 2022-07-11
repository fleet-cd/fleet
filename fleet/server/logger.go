package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Middleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path

		ctx.Next()

		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()
		timeStamp := time.Now()
		latency := timeStamp.Sub(start)

		logger.Info().
			Int("status_code", statusCode).
			Str("method", method).
			Str("path", path).
			Int64("latency_ms", latency.Milliseconds()).
			Send()

		ctx.Next()
	}
}
