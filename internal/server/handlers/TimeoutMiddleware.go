package handlers

import (
	"context"
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/config"
	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware is a Gin middleware function that sets a timeout for request execution.
// It creates a new context with a timeout based on the configured WriteTimeout from the HTTP server configuration.
// The middleware logs the timeout duration and ensures the request is canceled if it exceeds the timeout.
// The updated context is passed to the next handler in the chain.
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
