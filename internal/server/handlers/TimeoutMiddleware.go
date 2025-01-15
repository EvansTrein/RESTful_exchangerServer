package handlers

import (
	"context"
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/config"
	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(log *slog.Logger, conf *config.HTTPServer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "TimeoutMiddleware: call"

		newCtx, cancel := context.WithTimeout(ctx.Request.Context(), conf.WriteTimeout)
		defer cancel()

		ctx.Request = ctx.Request.WithContext(newCtx)

		log.Debug("an execution timeout has been set for the request",
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
			slog.String("request timeout", conf.WriteTimeout.String()),
		)

		ctx.Next()
	}
}
