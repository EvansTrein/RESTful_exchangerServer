package handlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/config"
	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(log *slog.Logger, conf *config.HTTPServer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "TimeoutMiddleware: call"

		newCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Second * 5)
		defer cancel()

		ctx.Request = ctx.Request.WithContext(newCtx)

		log.Debug("an execution timeout has been set for the request",
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method))

		ctx.Next()
	}
}
